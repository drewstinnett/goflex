package goflex

import (
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// PlaylistType describes what type a given playlist is.
type PlaylistType string

var (
	// VideoPlaylist is a playlist of just videos
	VideoPlaylist PlaylistType = "video"
	// AudioPlaylist is just audio
	AudioPlaylist PlaylistType = "audio"
	// PhotoPlaylist is just photos
	PhotoPlaylist PlaylistType = "photo"
)

// PlaylistService describes how a playlist service operates.
type PlaylistService interface {
	List() (map[PlaylistTitle]*Playlist, error)
	GetWithName(PlaylistTitle) (*Playlist, error)
	Create(PlaylistTitle, PlaylistType, bool) error
	GetOrCreate(PlaylistTitle, PlaylistType, bool) (*Playlist, bool, error)
	Delete(int) error
	DeleteEpisode(PlaylistTitle, ShowTitle, SeasonNumber, EpisodeNumber) error
	Exists(PlaylistTitle) (bool, error)
	Clear(int) error
	InsertEpisodes(int, EpisodeList) error
	Randomize(RandomizeRequest) (*RandomizeResponse, error)
	Episodes(Playlist) (EpisodeList, error)
	EpisodeID(Playlist, ShowTitle, SeasonNumber, EpisodeNumber) (int, error)
}

// PlaylistServiceOp is the operator for the PlaylistService.
type PlaylistServiceOp struct {
	p *Flex
}

// PlaylistEpisodeCache describes the playlist episodes
// type PlaylistEpisodeCache map[string]map[int]map[int]int

// PlaylistTitle is the title of a playlist.
type PlaylistTitle string

// Playlist is the important identifiers of a playlist.
type Playlist struct {
	ID       int
	Title    PlaylistTitle
	GUID     string
	Duration time.Duration
}

// RandomizeSeries defines parameters on randomizing a specific show
type RandomizeSeries struct {
	Filter       EpisodeFilter `json:"episodes" yaml:"episodes"`
	LookbackDays int           `json:"lookback" yaml:"lookback_days"`
}

// RandomizeRequestOpt defines how you request a new RandomizeRequest.
type RandomizeRequestOpt func(*RandomizeRequest)

// RandomizeRequestList is a list of RandomizeRequests.
type RandomizeRequestList []RandomizeRequest

// RandomizeRequest decides how to refill a Playlist.
type RandomizeRequest struct {
	Playlist PlaylistTitle     `yaml:"playlist"`
	Series   []RandomizeSeries `yaml:"series"`
	RefillAt int               `yaml:"refill_at"`
}

// NewRandomizeRequest returns a new RandomizeRequest using functional options
func NewRandomizeRequest(
	playlist PlaylistTitle,
	series []RandomizeSeries,
	opts ...RandomizeRequestOpt,
) (*RandomizeRequest, error) {
	req := RandomizeRequest{
		Playlist: playlist,
		Series:   series,
		RefillAt: 5,
	}
	for _, opt := range opts {
		opt(&req)
	}
	if req.Playlist == "" {
		return nil, errors.New("playlist must not be empty")
	}
	if len(req.Series) == 0 {
		return nil, errors.New("series muset not be empty")
	}
	return &req, nil
}

// RandomizeResponse is what we get back from requesting a Playlist be randomized.
type RandomizeResponse struct {
	// Refilled         bool          `json:"refilled,omitempty"`
	RefillReason     string        `json:"reason,omitempty"`
	Created          bool          `json:"created,omitempty"`
	Removed          EpisodeList   `json:"removed,omitempty"`
	Remaining        EpisodeList   `json:"remaining,omitempty"`
	OriginalEpisodes EpisodeList   `json:"original_episodes,omitempty"`
	UnviewedEpisodes EpisodeList   `json:"unviewed_episodes,omitempty"`
	SleepFor         time.Duration `json:"next_check,omitempty"`
}

func (svc *PlaylistServiceOp) processCreation(resp *RandomizeResponse, playlist *Playlist) error {
	// If created, always do a refill
	if resp.Created {
		resp.RefillReason = "newly created playlist"
	} else {
		// Otherwise get the original episodes
		var err error
		resp.OriginalEpisodes, err = svc.Episodes(*playlist)
		if err != nil {
			return err
		}
		slog.Debug("found existing playlist", "count", len(resp.OriginalEpisodes))
	}
	return nil
}

func (svc *PlaylistServiceOp) processViewed(
	resp *RandomizeResponse,
	req RandomizeRequest,
) (map[ShowTitle]EpisodeList, error) {
	viewedMap := make(map[ShowTitle]EpisodeList, len(req.Series))

	// collect the total removed and remaining in all of the series below
	for _, series := range req.Series {
		exists, err := svc.p.Shows.Exists(series.Filter.Show)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errors.New("show does not exist: " + string(series.Filter.Show))
		}

		// Get viewed
		since := -daysToDuration(series.LookbackDays)
		slog.Debug("looking back for episodes", "days", series.LookbackDays, "since", since, "show", series.Filter.Show)
		viewedMap[series.Filter.Show], err = svc.p.Sessions.HistoryEpisodes(
			time.Now().Add(since),
			series.Filter.Show,
		)
		if err != nil {
			return nil, err
		}
		slog.Debug("found viewed episodes", "count", len(viewedMap[series.Filter.Show]))

		// Figure out remaining
		remaining, removed := resp.OriginalEpisodes.Subtract(viewedMap[series.Filter.Show])
		slog.Debug(
			"removed viewed episodes",
			"removed",
			len(removed),
			"remaining",
			len(remaining),
			"show",
			series.Filter.Show,
		)
		resp.Remaining = append(resp.Remaining, remaining...)
		resp.Removed = append(resp.Removed, removed...)
	}
	return viewedMap, nil
}

func (svc *PlaylistServiceOp) initRandomize(req RandomizeRequest) (*RandomizeResponse, *Playlist, error) {
	resp := &RandomizeResponse{}
	// Inspect the playlist, create it if it doesn't exist
	playlist, created, err := svc.p.Playlists.GetOrCreate(req.Playlist, VideoPlaylist, false)
	if err != nil {
		return nil, nil, err
	}
	resp.Created = created
	return resp, playlist, nil
}

// Randomize randomizes a playlist with episodes from given series.
func (svc *PlaylistServiceOp) Randomize(req RandomizeRequest) (*RandomizeResponse, error) {
	resp, playlist, err := svc.initRandomize(req)
	if err != nil {
		return nil, err
	}

	if err := svc.processCreation(resp, playlist); err != nil {
		return nil, err
	}

	viewedMap, err := svc.processViewed(resp, req)
	if err != nil {
		return nil, err
	}

	// Did we dip below the refill line?
	// DREW: Figure out the runtime and base refill on that
	// slog.Info("remain", "remaining", resp.Remaining.Runtime())
	if (resp.RefillReason == "") && len(resp.Remaining) <= req.RefillAt {
		resp.RefillReason = fmt.Sprintf("playlist dipped below %v, was at: %v", req.RefillAt, len(resp.Remaining))
	}

	// Remove all the stuff we have already seen
	if err := svc.removeSeen(resp, req, *playlist); err != nil {
		return nil, err
	}

	// Refill if necessary
	if resp.RefillReason != "" {
		if err := svc.refillRand(resp, req, *playlist, viewedMap); err != nil {
			return nil, err
		}
	}

	// Figure out when we should check again
	if resp.SleepFor, err = svc.sleepFor(*playlist); err != nil {
		return nil, err
	}
	return resp, nil
}

func (svc *PlaylistServiceOp) sleepFor(playlist Playlist) (time.Duration, error) {
	episodes, err := svc.Episodes(playlist)
	if err != nil {
		return svc.p.maxSleep, err
	}
	if len(episodes) > 0 {
		nextEpisode, err := svc.p.Shows.Episode(episodes[0].Show, episodes[0].Season, episodes[0].Episode)
		if err != nil {
			return svc.p.maxSleep, err
		}
		slog.Info("Next", "episode", nextEpisode)
		currentlyWatching, err := svc.p.Sessions.ActiveEpisodes()
		if err != nil {
			return svc.p.maxSleep, err
		}
		for _, item := range currentlyWatching {
			slog.Info("looking at", "show", item.Show, "next", nextEpisode.Show, "duration", item.Duration)
			if item.Show == nextEpisode.Show {
				slog.Debug("found currently watching episode", "episode", item.String(), "duration", item.Duration)
				return item.Duration, nil
			}
		}
	}
	return svc.p.maxSleep, nil
}

func (svc *PlaylistServiceOp) refillRand(
	resp *RandomizeResponse,
	req RandomizeRequest,
	playlist Playlist,
	viewedMap map[ShowTitle]EpisodeList,
) error {
	slog.Debug("attempting to refill playlist", "playlist", req.Playlist, "reason", resp.RefillReason)
	if err := svc.Clear(playlist.ID); err != nil {
		return err
	}
	for _, series := range req.Series {
		shows, err := svc.p.Shows.Match(series.Filter.Show)
		if err != nil {
			return err
		}

		// allEpisodes, err := shows.EpisodesWithFilter(EpisodeFilter{
		allEpisodes, err := svc.p.Shows.EpisodesWithFilter(shows, EpisodeFilter{
			LatestSeason:   series.Filter.LatestSeason,
			EarliestSeason: series.Filter.EarliestSeason,
		})
		if err != nil {
			return err
		}

		unviewedEpisodes, _ := allEpisodes.Subtract(viewedMap[series.Filter.Show])
		rand.Shuffle(len(unviewedEpisodes), func(i, j int) {
			unviewedEpisodes[i], unviewedEpisodes[j] = unviewedEpisodes[j], unviewedEpisodes[i]
		})
		resp.UnviewedEpisodes = append(resp.UnviewedEpisodes, unviewedEpisodes...)
	}

	if len(resp.UnviewedEpisodes) < req.RefillAt {
		return fmt.Errorf(
			"not enough unwatched episodes to refill. unwatched: %v, refill-at: %v",
			len(resp.UnviewedEpisodes),
			req.RefillAt,
		)
	} else {
		slog.Info("refilling playlist", "title", playlist.Title, "episodes", len(resp.UnviewedEpisodes), "reason", resp.RefillReason)
		return svc.InsertEpisodes(playlist.ID, resp.UnviewedEpisodes)
	}
}

func (svc *PlaylistServiceOp) removeSeen(resp *RandomizeResponse, req RandomizeRequest, playlist Playlist) error {
	// Remove things we have seen
	if len(resp.Removed) > 0 {
		slog.Debug(
			"New length of episodes after removing viewed",
			"remaining",
			len(resp.Remaining),
			"removed",
			len(resp.Removed),
			"original",
			len(resp.OriginalEpisodes),
		)
		for _, item := range resp.Removed {
			slog.Info("removing episode", "playlist", req.Playlist, "episode", item.String())
			if err := svc.DeleteEpisode(playlist.Title, item.Show, item.Season, item.Episode); err != nil {
				return err
			}
		}
	}
	return nil
}

// InsertEpisodes inserts an episode in to a playlist.
func (svc *PlaylistServiceOp) InsertEpisodes(playlistID int, episodes EpisodeList) error {
	if len(episodes) == 0 {
		return nil
	}
	ids := make([]string, len(episodes))
	for idx, item := range episodes {
		ids[idx] = fmt.Sprint(item.ID)
	}
	machineID, err := svc.p.Server.MachineID()
	if err != nil {
		return err
	}
	var ret struct{}
	if err := svc.p.sendRequestXML(mustNewRequest("PUT",
		fmt.Sprintf("%v/playlists/%v/items?uri=%v",
			svc.p.baseURL,
			playlistID,
			fmt.Sprintf("server://%v/com.plexapp.plugins.library/library/metadata/%v", machineID, strings.Join(ids, ",")),
		)),
		&ret, nil); err != nil {
		return err
	}
	return nil
}

// Clear removes all items from a playlist.
func (svc *PlaylistServiceOp) Clear(id int) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%v/playlists/%v/items", svc.p.baseURL, id), nil)
	if err != nil {
		return err
	}
	var res struct{}
	if err := svc.p.sendRequestXML(req, &res, nil); err != nil {
		return err
	}
	// TODO: Clear the playlist cache here
	svc.p.cache.DeletePrefix("playlist-episodes-" + fmt.Sprint(id))
	return nil
}

// Delete deletees a playlist.
func (svc *PlaylistServiceOp) Delete(id int) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%v/playlists/%v", svc.p.baseURL, id), nil)
	if err != nil {
		return err
	}
	var res struct{}
	if err := svc.p.sendRequestXML(req, &res, nil); err != nil {
		return err
	}
	// DREW: Clear the playlist cache here
	return nil
}

// GetOrCreate returns a playlist with the given title, or creates a new one with given kind and smart options.
func (svc *PlaylistServiceOp) GetOrCreate(title PlaylistTitle, kind PlaylistType, smart bool) (*Playlist, bool, error) {
	var created bool
	exists, err := svc.Exists(title)
	if err != nil {
		return nil, false, fmt.Errorf("error checking if exists in GetOrCreate: %w", err)
	}
	if !exists {
		if err := svc.Create(title, kind, smart); err != nil {
			return nil, false, fmt.Errorf("error creating playlist in GetOrCreate: %w", err)
		}
		created = true
	}
	got, err := svc.GetWithName(title)
	if err != nil {
		return nil, created, fmt.Errorf("error getting with name in GetOrCreate: %w", err)
	}
	return got, created, nil
}

// Exists returns true if a playlist already exists.
func (svc *PlaylistServiceOp) Exists(n PlaylistTitle) (bool, error) {
	items, err := svc.List()
	if err != nil {
		return false, err
	}
	_, ok := items[n]
	return ok, nil
}

// Create creates a new playlist.
func (svc *PlaylistServiceOp) Create(title PlaylistTitle, kind PlaylistType, smart bool) error {
	exists, err := svc.Exists(title)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("playlist with title already exists: " + string(title))
	}
	smartInt := 0
	if smart {
		smartInt = 1
	}
	machineID, err := svc.p.Server.MachineID()
	if err != nil {
		return err
	}
	req := mustNewRequest(
		"POST",
		fmt.Sprintf(
			"%v/playlists?type=%v&title=%v&smart=%v&uri=server://%v/com.plexapp.plugins.library/",
			svc.p.baseURL,
			kind,
			url.QueryEscape(string(title)),
			smartInt,
			machineID,
		),
	)
	var res CreatePlaylistResponse
	if err := svc.p.sendRequestXML(req, &res, nil); err != nil {
		return err
	}
	svc.p.cache.DeletePrefix("playlists")
	return nil
}

// List lists out playlists on a plex server.
func (svc *PlaylistServiceOp) List() (map[PlaylistTitle]*Playlist, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%v/playlists", svc.p.baseURL), nil)
	if err != nil {
		return nil, err
	}
	var pr PlaylistsResponse
	if err := svc.p.sendRequestXML(req, &pr, &cacheConfig{prefix: "playlists", ttl: time.Minute * 10}); err != nil {
		return nil, err
	}

	ret := map[PlaylistTitle]*Playlist{}
	for _, item := range pr.Playlist {
		id, err := strconv.Atoi(item.RatingKey)
		if err != nil {
			return nil, err
		}
		second := 0
		if item.Duration != "" {
			var err error
			if second, err = strconv.Atoi(item.Duration); err != nil {
				return nil, err
			}
		}
		ret[PlaylistTitle(item.Title)] = &Playlist{
			ID:       id,
			Title:    PlaylistTitle(item.Title),
			GUID:     item.GUID,
			Duration: time.Duration(second) * time.Millisecond,
		}
	}
	return ret, nil
}

// GetWithName returns a playlist by name.
func (svc *PlaylistServiceOp) GetWithName(n PlaylistTitle) (*Playlist, error) {
	items, err := svc.List()
	if err != nil {
		return nil, err
	}
	got, ok := items[n]
	if !ok {
		return nil, errors.New("playlist not found with name: " + string(n))
	}
	return got, nil
}

func episodeWithVideo(item Video) (*Episode, error) {
	parentI, err := strconv.Atoi(item.ParentIndex)
	if err != nil {
		return nil, err
	}
	index, err := strconv.Atoi(item.Index)
	if err != nil {
		return nil, err
	}

	// Figure out how long the episode is.
	du, err := strconv.Atoi(item.Duration)
	if err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(item.RatingKey)
	if err != nil {
		return nil, err
	}

	var viewedInt int64
	if item.LastViewedAt != "" {
		var err error
		viewedInt, err = strconv.ParseInt(item.LastViewedAt, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	playlistID, err := strconv.Atoi(item.PlaylistItemID)
	if err != nil {
		return nil, err
	}
	vc := 0
	if item.ViewCount != "" {
		var err error
		if vc, err = strconv.Atoi(item.ViewCount); err != nil {
			return nil, err
		}
	}
	viewed := time.Unix(viewedInt, 0)
	e := &Episode{
		ID:             id,
		PlaylistItemID: playlistID,
		Title:          item.Title,
		Show:           ShowTitle(item.GrandparentTitle),
		Season:         SeasonNumber(parentI),
		Episode:        EpisodeNumber(index),
		Watched:        &viewed,
		Duration:       time.Duration(du) * time.Millisecond,
		ViewCount:      vc,
	}
	if item.LastViewedAt != "" {
		var viewedInt int64
		var err error
		viewedInt, err = strconv.ParseInt(item.LastViewedAt, 10, 64)
		if err != nil {
			return nil, err
		}
		viewed := time.Unix(viewedInt, 0)
		e.Watched = &viewed
	}
	return e, nil
}

// Episodes returns a new playlist by the name.
func (svc *PlaylistServiceOp) Episodes(p Playlist) (EpisodeList, error) {
	if p.Title == "" {
		return nil, errors.New("playlist Title must not be empty")
	}
	ret := EpisodeList{}
	var plr PlaylistResponse
	if err := svc.p.sendRequestXML(mustNewRequest(http.MethodGet, fmt.Sprintf("%v/playlists/%v/items", svc.p.baseURL, p.ID)), &plr,
		&cacheConfig{prefix: "playlist-episodes-" + fmt.Sprint(p.ID), ttl: time.Minute * 60}); err != nil {
		return nil, err
	}
	for _, item := range plr.Video {
		episode, err := episodeWithVideo(item)
		if err != nil {
			return nil, err
		}
		ret = append(ret, *episode)
	}
	return ret, nil
}

// DeleteItem removes an item from the given playlist.
func (svc *PlaylistServiceOp) DeleteItem(id int, keys ...int) error {
	for _, k := range keys {
		req, err := http.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("%v/playlists/%v/items/%v", svc.p.baseURL, id, k),
			nil,
		)
		if err != nil {
			return err
		}
		var got struct{}
		if err := svc.p.sendRequestXML(req, &got, nil); err != nil {
			return err
		}
	}
	return nil
}

func (svc *PlaylistServiceOp) EpisodeID(
	playlist Playlist,
	show ShowTitle,
	season SeasonNumber,
	episode EpisodeNumber,
) (int, error) {
	if show == "" {
		return 0, errors.New("must specify a show title")
	}
	if season == 0 {
		return 0, errors.New("must specify a season")
	}
	if episode == 0 {
		return 0, errors.New("must specify an episode")
	}
	episodes, err := svc.Episodes(playlist)
	if err != nil {
		return 0, err
	}
	for _, episodeI := range episodes {
		if (show == episodeI.Show) && (season == episodeI.Season) && (episode == episodeI.Episode) {
			return episodeI.PlaylistItemID, nil
		}
	}
	return 0, errors.New("episode not found")
}

// DeleteEpisode removes an item by title, season number,  episode number.
func (svc *PlaylistServiceOp) DeleteEpisode(
	playlist PlaylistTitle,
	show ShowTitle,
	season SeasonNumber,
	episode EpisodeNumber,
) error {
	if show == "" {
		return errors.New("cannot delete episode with empty show")
	}
	pl, err := svc.GetWithName(playlist)
	if err != nil {
		return err
	}

	k, err := svc.EpisodeID(*pl, show, season, episode)
	if err != nil {
		return err
	}
	return svc.DeleteItem(pl.ID, k)
}
