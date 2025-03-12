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

	"github.com/drewstinnett/gout/v2"
)

// PlaylistType describes what type a given playlist is
type PlaylistType string

var (
	// VideoPlaylist is a playlist of just videos
	VideoPlaylist PlaylistType = "video"
	// AudioPlaylist is just audio
	AudioPlaylist PlaylistType = "audio"
	// PhotoPlaylist is just photos
	PhotoPlaylist PlaylistType = "photo"
)

// PlaylistService describes how a playlist service operates
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

// PlaylistServiceOp is the operator for the PlaylistService
type PlaylistServiceOp struct {
	p *Plex
	// cache map[PlaylistTitle]*Playlist
}

// PlaylistEpisodeCache describes the playlist episodes
type PlaylistEpisodeCache map[string]map[int]map[int]int

// PlaylistTitle is the title of a playlist
type PlaylistTitle string

// Playlist is the important identifiers of a playlist
type Playlist struct {
	ID       int
	Title    PlaylistTitle
	GUID     string
	Duration time.Duration
}

// RandomizeSeries defines parameters on randomizing a specific show
type RandomizeSeries struct {
	Filter   EpisodeFilter  `json:"episodes"`
	Lookback *time.Duration `json:"lookback"`
	RefillAt int            `json:"refill_at"`
}

// RandomizeRequestOpt defines how you request a new RandomizeRequest
type RandomizeRequestOpt func(*RandomizeRequest)

// RandomizeRequest decides how to refill a Playlist
type RandomizeRequest struct {
	playlist PlaylistTitle
	series   []RandomizeSeries
	refillAt int
}

// / NewRandomizeRequest returns a new RandomizeRequest using functional options
func NewRandomizeRequest(playlist PlaylistTitle, series []RandomizeSeries, opts ...RandomizeRequestOpt) (*RandomizeRequest, error) {
	req := RandomizeRequest{
		playlist: playlist,
		series:   series,
		refillAt: 5,
	}
	for _, opt := range opts {
		opt(&req)
	}
	if req.playlist == "" {
		return nil, errors.New("playlist must not be empty")
	}
	if len(req.series) == 0 {
		return nil, errors.New("series muset not be empty")
	}
	return &req, nil
}

// RandomizeResponse is what we get back from requesting a Playlist be randomized
type RandomizeResponse struct {
	// Refilled         bool          `json:"refilled,omitempty"`
	RefillReason     string        `json:"reason,omitempty"`
	Created          bool          `json:"created,omitempty"`
	Removed          EpisodeList   `json:"removed,omitempty"`
	Remaining        EpisodeList   `json:"remaining,omitempty"`
	OriginalEpisodes EpisodeList   `json:"original_episodes,omitempty"`
	UnviewedEpisodes EpisodeList   `json:"unviewed_episodes,omitempty"`
	NextCheck        time.Duration `json:"next_check,omitempty"`
}

// Randomize randomizes a playlist with episodes from given series
func (svc *PlaylistServiceOp) Randomize(req RandomizeRequest) (*RandomizeResponse, error) {
	resp := &RandomizeResponse{}

	// Inspect the playlist, create it if it doesn't exist
	playlist, created, err := svc.p.Playlists.GetOrCreate(req.playlist, VideoPlaylist, false)
	if err != nil {
		return nil, err
	}
	resp.Created = created

	// If created, always do a refill
	if resp.Created {
		resp.RefillReason = "newly created playlist"
	} else {
		// Otherwise get the original episodes
		var err error
		// resp.OriginalEpisodes, err = playlist.Episodes()
		resp.OriginalEpisodes, err = svc.Episodes(*playlist)
		if err != nil {
			return nil, err
		}
		slog.Debug("found existing playlist", "count", len(resp.OriginalEpisodes))
	}

	viewedMap := make(map[ShowTitle]EpisodeList, len(req.series))

	// collect the total removed and remaining in all of the series below
	for _, series := range req.series {
		exists, err := svc.p.Shows.Exists(series.Filter.Show)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errors.New("show does not exist: " + string(series.Filter.Show))
		}

		// Get viewed
		viewedMap[series.Filter.Show], err = svc.p.Sessions.HistoryEpisodes(toPTR(time.Now().Add(-fromPTR(series.Lookback))), series.Filter.Show)
		if err != nil {
			return nil, err
		}
		slog.Debug("found viewed episodes", "count", len(viewedMap[series.Filter.Show]))

		// Figure out remaining
		remaining, removed := resp.OriginalEpisodes.Subtract(viewedMap[series.Filter.Show])
		slog.Debug("removed viewed episodes", "removed", len(removed), "remaining", len(remaining), "show", series.Filter.Show)
		resp.Remaining = append(resp.Remaining, remaining...)
		resp.Removed = append(resp.Removed, removed...)

	}

	// Did we dip below the refill line?
	if (resp.RefillReason == "") && len(resp.Remaining) <= req.refillAt {
		resp.RefillReason = fmt.Sprintf("playlist dipped below %v, was at: %v", req.refillAt, len(resp.Remaining))
	}

	// Remove things we have seen
	if len(resp.Removed) > 0 {
		slog.Debug("New length of episodes after removing viewed", "remaining", len(resp.Remaining), "removed", len(resp.Removed), "original", len(resp.OriginalEpisodes))
		for _, item := range resp.Removed {
			slog.Info("removing episode", "playlist", req.playlist, "episode", item.String())
			if err := svc.DeleteEpisode(playlist.Title, item.Show, item.Season, item.Episode); err != nil {
				return nil, err
			}
		}
	}

	if resp.RefillReason != "" {
		slog.Debug("attempting to refill playlist", "playlist", req.playlist, "reason", resp.RefillReason)
		if err := svc.Clear(playlist.ID); err != nil {
			return nil, err
		}

		for _, series := range req.series {
			shows, err := svc.p.Shows.Match(series.Filter.Show)
			if err != nil {
				return nil, err
			}
			gout.MustPrint(shows)

			// allEpisodes, err := shows.EpisodesWithFilter(EpisodeFilter{
			allEpisodes, err := svc.p.Shows.EpisodesWithFilter(shows, EpisodeFilter{
				LatestSeason:   series.Filter.LatestSeason,
				EarliestSeason: series.Filter.EarliestSeason,
			})
			if err != nil {
				return nil, err
			}

			unviewedEpisodes, _ := allEpisodes.Subtract(viewedMap[series.Filter.Show])
			rand.Shuffle(len(unviewedEpisodes), func(i, j int) {
				unviewedEpisodes[i], unviewedEpisodes[j] = unviewedEpisodes[j], unviewedEpisodes[i]
			})
			resp.UnviewedEpisodes = append(resp.UnviewedEpisodes, unviewedEpisodes...)
		}

		if len(resp.UnviewedEpisodes) < req.refillAt {
			return nil, fmt.Errorf("not enough unwatched episodes to refill. unwatched: %v, refill-at: %v", len(resp.UnviewedEpisodes), req.refillAt)
		} else {
			slog.Info("refilling playlist", "title", playlist.Title, "episodes", len(resp.UnviewedEpisodes), "reason", resp.RefillReason)
			return resp, svc.InsertEpisodes(playlist.ID, resp.UnviewedEpisodes)
		}
	}
	return resp, nil
}

// InsertEpisodes inserts an episode in to a playlist
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
			// episode.URI(),
		)),
		&ret, nil); err != nil {
		return err
	}
	return nil
}

// Clear removes all items from a playlist
func (svc *PlaylistServiceOp) Clear(id int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%v/playlists/%v/items", svc.p.baseURL, id), nil)
	if err != nil {
		return err
	}
	var res struct{}
	if err := svc.p.sendRequestXML(req, &res, nil); err != nil {
		return err
	}
	// TODO: Clear the playlist cache here
	return nil
}

// Delete deletees a playlist
func (svc *PlaylistServiceOp) Delete(id int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%v/playlists/%v", svc.p.baseURL, id), nil)
	if err != nil {
		return err
	}
	var res struct{}
	if err := svc.p.sendRequestXML(req, &res, nil); err != nil {
		return err
	}
	// TODO: Clear the playlist cache here
	return nil
}

// GetOrCreate returns a playlist with the given title, or creates a new one with given kind and smart options
func (svc *PlaylistServiceOp) GetOrCreate(title PlaylistTitle, kind PlaylistType, smart bool) (*Playlist, bool, error) {
	var created bool
	exists, err := svc.Exists(title)
	if err != nil {
		return nil, false, err
	}
	if !exists {
		if err := svc.Create(title, kind, smart); err != nil {
			return nil, false, err
		}
		created = true
	}
	got, err := svc.GetWithName(title)
	if err != nil {
		return nil, created, err
	}
	return got, created, nil
}

// Exists returns true if a playlist already exists
func (svc *PlaylistServiceOp) Exists(n PlaylistTitle) (bool, error) {
	items, err := svc.List()
	if err != nil {
		return false, err
	}
	_, ok := items[n]
	return ok, nil
}

// Create creates a new playlist
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
	req := mustNewRequest("POST", fmt.Sprintf("%v/playlists?type=%v&title=%v&smart=%v&uri=server://%v/com.plexapp.plugins.library/", svc.p.baseURL, kind, url.QueryEscape(string(title)), smartInt, machineID))
	var res CreatePlaylistResponse
	if err := svc.p.sendRequestXML(req, &res, nil); err != nil {
		return err
	}
	// TODO: Clear playlist cache
	return nil
}

// List lists out playlists on a plex server
func (svc *PlaylistServiceOp) List() (map[PlaylistTitle]*Playlist, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/playlists", svc.p.baseURL), nil)
	if err != nil {
		return nil, err
	}
	var pr PlaylistsResponse
	if err := svc.p.sendRequestXML(req, &pr, toPTR(time.Minute*10)); err != nil {
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
			Duration: time.Duration(second) * time.Second,
		}
	}
	return ret, nil
	/*
		if svc.cache == nil {
			if err := svc.updateCache(); err != nil {
				return nil, err
			}
		}
		return svc.cache, nil
	*/
}

// GetWithName returns a playlist by name
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

/*
func (svc *PlaylistServiceOp) updateCache() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/playlists", svc.p.baseURL), nil)
	if err != nil {
		return err
	}
	var pr PlaylistsResponse
	if err := svc.p.sendRequestXML(req, &pr, nil); err != nil {
		return err
	}

	svc.cache = map[PlaylistTitle]*Playlist{}
	for _, item := range pr.Playlist {
		id, err := strconv.Atoi(item.RatingKey)
		if err != nil {
			return err
		}
		second := 0
		if item.Duration != "" {
			var err error
			if second, err = strconv.Atoi(item.Duration); err != nil {
				return err
			}
		}
		svc.cache[PlaylistTitle(item.Title)] = &Playlist{
			ID:       id,
			Title:    PlaylistTitle(item.Title),
			GUID:     item.GUID,
			Duration: time.Duration(second) * time.Second,
			// p:        svc.p,
		}
	}
	return nil
}
*/

// Episodes returns a new playlist by the name
func (svc *PlaylistServiceOp) Episodes(p Playlist) (EpisodeList, error) {
	if p.Title == "" {
		return nil, errors.New("playlist Title must not be empty")
	}
	ret := EpisodeList{}
	var plr PlaylistResponse
	if err := svc.p.sendRequestXML(mustNewRequest("GET", fmt.Sprintf("%v/playlists/%v/items", svc.p.baseURL, p.ID)), &plr, toPTR(time.Minute*60)); err != nil {
		return nil, err
	}
	for _, item := range plr.Video {
		parentI, err := strconv.Atoi(item.ParentIndex)
		if err != nil {
			return nil, err
		}
		index, err := strconv.Atoi(item.Index)
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
		ret = append(ret, Episode{
			ID:             id,
			PlaylistItemID: playlistID,
			Title:          item.Title,
			Show:           ShowTitle(item.GrandparentTitle),
			Season:         SeasonNumber(parentI),
			Episode:        EpisodeNumber(index),
			Watched:        &viewed,
			ViewCount:      vc,
		})
	}
	return ret, nil
}

/*
// Episodes returns a new playlist by the name
func (l *Playlist) Episodes() (EpisodeList, error) {
	var plr PlaylistResponse
	if err := l.p.sendRequestXML(mustNewRequest("GET", fmt.Sprintf("%v/playlists/%v/items", l.p.baseURL, l.ID)), &plr); err != nil {
		return nil, err
	}
	ret := EpisodeList{}
	for _, item := range plr.Video {
		parentI, err := strconv.Atoi(item.ParentIndex)
		if err != nil {
			return nil, err
		}
		index, err := strconv.Atoi(item.Index)
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
		ret = append(ret, Episode{
			ID:             id,
			PlaylistItemID: playlistID,
			Title:          item.Title,
			Show:           item.GrandparentTitle,
			Season:         parentI,
			Episode:        index,
			Watched:        &viewed,
			ViewCount:      vc,
			p:              l.p,
		})
	}

	return ret, nil
}
*/

// DeleteItem removes an item from the given playlist
func (svc *PlaylistServiceOp) DeleteItem(id int, keys ...int) error {
	for _, k := range keys {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%v/playlists/%v/items/%v", svc.p.baseURL, id, k), nil)
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

func (svc *PlaylistServiceOp) EpisodeID(playlist Playlist, show ShowTitle, season SeasonNumber, episode EpisodeNumber) (int, error) {
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
	return 0, errors.New("episode now found")
}

/*
func (l Playlist) episodeID(show string, season, episode int) (int, error) {
	if show == "" {
		return 0, errors.New("must specify a show title")
	}
	if season == 0 {
		return 0, errors.New("must specify a season")
	}
	if episode == 0 {
		return 0, errors.New("must specify an episode")
	}
	episodes, err := l.p.Playlists.Episodes(l)
	if err != nil {
		return 0, err
	}
	for _, episodeI := range episodes {
		if (show == episodeI.Show) && (season == episodeI.Season) && (episode == episodeI.Episode) {
			return episodeI.PlaylistItemID, nil
		}
	}
	return 0, errors.New("episode now found")
}
*/

// DeleteEpisode removes an item by title, season number,  episode number
func (svc *PlaylistServiceOp) DeleteEpisode(playlist PlaylistTitle, show ShowTitle, season SeasonNumber, episode EpisodeNumber) error {
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

/*
func (svc *PlaylistServiceOp) updateEpisodeCache(episodes *EpisodeList) {
	c := PlaylistEpisodeCache{}
	for _, e := range *episodes {
		if _, ok := c[e.Show]; !ok {
			c[e.Show] = map[int]map[int]int{}
		}
		if _, ok := c[e.Show][e.Season]; !ok {
			c[e.Show][e.Season] = map[int]int{}
		}
		c[e.Show][e.Season][e.Episode] = e.PlaylistItemID
	}
	svc.episodeCache = c
}

func (l *Playlist) updateEpisodeCache(episodes *EpisodeList) {
	c := PlaylistEpisodeCache{}
	for _, e := range *episodes {
		if _, ok := c[e.Show]; !ok {
			c[e.Show] = map[int]map[int]int{}
		}
		if _, ok := c[e.Show][e.Season]; !ok {
			c[e.Show][e.Season] = map[int]int{}
		}
		c[e.Show][e.Season][e.Episode] = e.PlaylistItemID
	}
	l.episodeCache = c
}
*/
