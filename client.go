/*
Package plexrando does the randomization bits
*/
package plexrando

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/LukeHagar/plexgo"
	"github.com/LukeHagar/plexgo/models/operations"
	"moul.io/http2curl/v2"
)

const version string = "0.1.0"

// Plex connects to our custom stuff
type Plex struct {
	serverID       string
	baseURL        string
	token          string
	printCurl      bool
	client         *http.Client
	API            *plexgo.PlexAPI
	libraryMap     map[string]int
	playlistMap    map[string]Playlist
	showLibraryMap map[string]string
	episodeIndex   map[string]map[int]map[int]string
}

// New uses functional options for a new plex
func New(opts ...func(*Plex)) (*Plex, error) {
	p := &Plex{
		client: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(p)
	}
	if p.baseURL == "" {
		return nil, errors.New("must set plex baseurl")
	}
	if p.token == "" {
		return nil, errors.New("must set token")
	}

	// Set up a mainstream API attribute
	p.API = plexgo.New(
		plexgo.WithSecurity(p.token),
		plexgo.WithServerURL(p.baseURL),
		plexgo.WithClientID("313FF6D7-5795-45E3-874F-B8FCBFD5E587"),
		plexgo.WithClientName("plex-trueget"),
		plexgo.WithClientVersion(version),
	)
	if err := p.init(); err != nil {
		panic(err)
	}
	return p, nil
}

// WithHTTPClient sets the http client on a new Plex
func WithHTTPClient(c *http.Client) func(*Plex) {
	return func(p *Plex) {
		p.client = c
	}
}

// WithBaseURL sets the base url for a plex client
func WithBaseURL(s string) func(*Plex) {
	return func(p *Plex) {
		p.baseURL = s
	}
}

// WithToken sets the plex token for a client
func WithToken(s string) func(*Plex) {
	return func(p *Plex) {
		p.token = s
	}
}

// WithPrintCurl sets the curl debug printer to on
func WithPrintCurl() func(*Plex) {
	return func(p *Plex) {
		p.printCurl = true
	}
}

// Episode represents an episode of television
type Episode struct {
	ID             string // ID is pretty much just the RatingKey
	PlaylistItemID string
	Title          string
	Show           string
	Season         int
	Episode        int
	Watched        *time.Time
	p              *Plex
}

// EpisodeList is multiple Episodes
type EpisodeList []Episode

// Runtime returns the total runtime of all episodes
func (l EpisodeList) Runtime() time.Duration {
	return 0
}

// Seasons returns episodes matching a season start and stop
func (l EpisodeList) Seasons(start, end int) EpisodeList {
	ret := EpisodeList{}
	// If no start/end specified, return everything
	if (start == 0) && (end == 0) {
		return l
	}
	for _, ep := range l {
		switch {
		// If we just have a start
		case ((start > 0) && (end == 0)) && (ep.Season >= start):
			ret = append(ret, ep)
		// If we just have an end
		case ((end > 0) && (start == 0)) && (ep.Season <= end):
			ret = append(ret, ep)
		// If we have a start and and end
		case (ep.Season >= start) && (ep.Season <= end):
			ret = append(ret, ep)
		}
	}
	return ret
}

// Subtract removes items from a list. Returns the edited list, and a list of episodes that were removed
func (l *EpisodeList) Subtract(s EpisodeList) (EpisodeList, EpisodeList) {
	r := EpisodeList{}
	removed := EpisodeList{}
	ids := s.IDs()
	for _, item := range *l {
		if !slices.Contains(ids, item.ID) {
			r = append(r, item)
		} else {
			removed = append(removed, item)
		}
	}
	return r, removed
}

// URIList returns a list of URIs for the episodes
func (l EpisodeList) URIList() []string {
	ret := make([]string, len(l))
	for idx, item := range l {
		ret[idx] = item.URI()
	}
	return ret
}

// IDs returns a list of IDs for the episodes
func (l EpisodeList) IDs() []string {
	ret := make([]string, len(l))
	for idx, item := range l {
		ret[idx] = item.ID
	}
	return ret
}

// NewEpisode returns a new Episode from a operations.GetPlaylistContentsMetadata
func (p *Plex) NewEpisode(item operations.GetPlaylistContentsMetadata) Episode {
	return Episode{
		ID:    *item.RatingKey,
		Title: *item.Title,
		p:     p,
	}
}

// NewEpisodeWithChildrenMeta uses operations.GetMetadataChildrenMetadata to create a new episode
func (p *Plex) NewEpisodeWithChildrenMeta(item operations.GetMetadataChildrenMetadata) Episode {
	return Episode{
		ID:    *item.RatingKey,
		Title: *item.Title,
		p:     p,
	}
}

// NewEpisodeWithSession uses session info to get episode info
func (p *Plex) NewEpisodeWithSession(item operations.GetSessionHistoryMetadata) (*Episode, error) {
	if item.RatingKey == nil {
		return nil, errors.New("missing rating key")
	}
	if item.Title == nil {
		return nil, errors.New("missing title")
	}
	watched := time.Unix(int64(*item.ViewedAt), 0)
	return &Episode{
		ID:      *item.RatingKey,
		Show:    *item.GrandparentTitle,
		Season:  *item.ParentIndex,
		Episode: *item.Index,
		Title:   *item.Title,
		Watched: &watched,
		p:       p,
	}, nil
}

// URI returns the URI for an episode. This is the format the Playlist stuff needs
func (e Episode) URI() string {
	return fmt.Sprintf("server://%v/com.plexapp.plugins.library/library/metadata/%v", e.p.serverID, e.ID)
}

// String fulfills the Stringer interface
func (e Episode) String() string {
	var ret string
	switch {
	case e.Show == "":
		ret = fmt.Sprintf("%v - %v", e.ID, e.Title)
	default:
		ret = fmt.Sprintf("%v - S%02dE%02d", e.Show, e.Season, e.Episode)
	}
	if e.Watched != nil {
		ret = fmt.Sprintf("%v (Watched: %v)", ret, e.Watched.Format("2006-01-02 15:04"))
	}
	return ret
}

// WithAPI sets the PlexAPI on a new plex
func WithAPI(a *plexgo.PlexAPI) func(*Plex) {
	return func(p *Plex) {
		p.API = a
	}
}

func (p *Plex) init() error {
	if err := p.serverInfo(); err != nil {
		return err
	}
	if err := p.updateSections(); err != nil {
		return err
	}
	return p.updatePlaylists()
}

func (p *Plex) serverInfo() error {
	res, err := p.API.Server.GetServerIdentity(context.Background())
	if err != nil {
		return err
	}
	p.serverID = *res.Object.MediaContainer.MachineIdentifier
	return nil
}

func (p *Plex) updatePlaylists() error {
	p.playlistMap = map[string]Playlist{}
	ctx := context.Background()
	lib, err := p.API.Playlists.GetPlaylists(ctx, operations.PlaylistTypeVideo.ToPointer(), operations.QueryParamSmartZero.ToPointer())
	if err != nil {
		return err
	}
	for _, playlist := range lib.Object.MediaContainer.Metadata {
		k, err := strconv.ParseFloat(*playlist.RatingKey, 64)
		if err != nil {
			return err
		}
		duration := time.Duration(0)
		if playlist.Duration != nil {
			duration = time.Duration(*playlist.Duration) * time.Second
		}
		p.playlistMap[*playlist.Title] = Playlist{
			ID:       k,
			URI:      *playlist.GUID,
			Duration: duration,
			p:        p,
		}
	}
	return nil
}

func (p *Plex) updateSections() error {
	p.libraryMap = map[string]int{}
	ctx := context.Background()
	lib, err := p.API.Library.GetAllLibraries(ctx)
	if err != nil {
		return err
	}
	for _, libd := range lib.Object.MediaContainer.Directory {
		id, err := strconv.Atoi(libd.Key)
		if err != nil {
			return err
		}
		p.libraryMap[libd.Title] = id
	}
	return nil
}

// GetOrCreatePlaylist creates an empty playlist using the given name. If a playlist
// with that name already exists, it will be emptied
func (p *Plex) GetOrCreatePlaylist(s string) (*Playlist, bool, error) {
	var playlist Playlist
	var ok bool
	if playlist, ok = p.playlistMap[s]; !ok {
		res, err := p.API.Playlists.CreatePlaylist(context.Background(), operations.CreatePlaylistRequest{
			Title: s,
			Type:  operations.CreatePlaylistQueryParamTypeVideo,
		})
		if err != nil {
			return nil, false, err
		}
		k, err := strconv.ParseFloat(*res.Object.MediaContainer.Metadata[0].RatingKey, 64)
		if err != nil {
			return nil, false, err
		}
		// duration := res.
		return &Playlist{
			ID:  k,
			URI: *res.Object.MediaContainer.Metadata[0].GUID,
			p:   p,
		}, true, nil
	}
	/*
		if _, err := p.API.Playlists.ClearPlaylistContents(context.Background(), playlist.ID); err != nil {
			return nil, false, err
		}
	*/
	return &playlist, false, nil
}

// Clear empties the playlist of all contents
func (l *Playlist) Clear() error {
	_, err := l.p.API.Playlists.ClearPlaylistContents(context.Background(), l.ID)
	return err
}

// Viewed returns a list of recently viewed episodes
func (p *Plex) Viewed(library string, since time.Time) (EpisodeList, error) {
	inProcess := []string{}
	lib := p.libraryMap[library]
	current, err := p.API.Sessions.GetSessions(
		context.Background(),
	)
	if err != nil {
		return nil, err
	}
	for _, item := range current.Object.MediaContainer.Metadata {
		inProcess = append(inProcess, *item.RatingKey)
	}
	res, err := p.API.Sessions.GetSessionHistory(
		context.Background(),
		plexgo.String("viewedAt:desc"),
		plexgo.Int64(1),
		&operations.QueryParamFilter{},
		plexgo.Int64(int64(lib)),
	)
	if err != nil {
		return nil, err
	}
	ret := EpisodeList{}
	for _, episode := range res.Object.MediaContainer.Metadata {
		if episode.RatingKey == nil {
			slog.Debug("missing Rating key", "episode", episode)
			continue
		}
		if slices.Contains(inProcess, *episode.RatingKey) {
			continue
		}
		if time.Unix(int64(*episode.ViewedAt), 0).Before(since) {
			break
		}
		ep, err := p.NewEpisodeWithSession(episode)
		if err != nil {
			slog.Warn("error creating episode", "error", err)
			continue
		}
		ret = append(ret, *ep)
	}
	return ret, nil
}

// Episodes returns a list of episodes for a show
func (p *Plex) Episodes(library, show string) (EpisodeList, error) {
	id, ok := p.libraryMap[library]
	if !ok {
		return nil, fmt.Errorf("unknown library: %v", library)
	}
	ret := EpisodeList{}
	res, err := p.API.Library.GetLibraryItems(context.Background(), operations.GetLibraryItemsRequest{
		Tag:         "all",
		SectionKey:  id,
		IncludeMeta: operations.GetLibraryItemsQueryParamIncludeMetaEnable.ToPointer(),
	})
	if err != nil {
		return nil, err
	}
	for _, item := range res.Object.MediaContainer.Metadata {
		if item.Title != show {
			continue
		}
		rk, err := strconv.ParseFloat(item.RatingKey, 64)
		if err != nil {
			return nil, err
		}
		children, err := p.API.Library.GetMetadataChildren(context.Background(), rk, plexgo.String("Stream"))
		if err != nil {
			return nil, err
		}

		for _, season := range children.Object.MediaContainer.Metadata {
			rk, err := strconv.ParseFloat(*season.RatingKey, 64)
			if err != nil {
				return nil, err
			}
			children, err := p.API.Library.GetMetadataChildren(context.Background(), rk, plexgo.String("Stream"))
			if err != nil {
				return nil, err
			}
			for _, e := range children.Object.MediaContainer.Metadata {
				ret = append(ret, p.NewEpisodeWithChildrenMeta(e))
			}
		}
	}
	return ret, nil
}

type playlistEpisodeCache map[string]map[int]map[int]string

// Playlist is the important identifiers of a playlist
type Playlist struct {
	ID           float64
	URI          string
	Duration     time.Duration
	p            *Plex
	episodeCache playlistEpisodeCache
}

// AddEpisodes adds an episode or more to a given playlist
func (l *Playlist) AddEpisodes(episodes EpisodeList) error {
	slog.Debug("Adding episodes back in", "count", len(episodes))
	ret := &EpisodeList{}
	for idx, episode := range episodes {
		i := float64(idx)
		if _, err := l.p.API.Playlists.AddPlaylistContents(context.Background(), l.ID, episode.URI(), &i); err != nil {
			return err
		}
		slog.Debug("Adding", "episode", episode.String())
	}
	l.updateEpisodeCache(ret)
	return l.p.updatePlaylists()
}

// AddEpisode adds an episode to a given playlist
func (l *Playlist) AddEpisode(episode Episode) error {
	slog.Debug("Adding episode back in", "count", episode)
	idx := float64(len(l.episodeCache))
	if _, err := l.p.API.Playlists.AddPlaylistContents(context.Background(), l.ID, episode.URI(), &idx); err != nil {
		return err
	}
	slog.Debug("Adding", "episode", episode.String())
	return l.p.updatePlaylists()
}

// Playlist returns a new playlist by the name

// Playlist returns a new playlist by the name
func (p *Plex) Playlist(s string) (*Playlist, error) {
	pl, ok := p.playlistMap[s]
	if !ok {
		return nil, fmt.Errorf("playlist not found. requested: %v, available: %v", s, p.playlistMap)
	}
	return &pl, nil
}

// DeleteEpisode removes an item by title, season number,  episode number
func (l *Playlist) DeleteEpisode(show string, season, episode int) error {
	if show == "" {
		return errors.New("cannot delete episode with empty show")
	}
	k, err := l.EpisodeKey(show, season, episode)
	if err != nil {
		return err
	}
	return l.DeleteItem(k)
}

func (l *Playlist) updateEpisodeCache(episodes *EpisodeList) {
	c := playlistEpisodeCache{}
	for _, e := range *episodes {
		if _, ok := c[e.Show]; !ok {
			c[e.Show] = map[int]map[int]string{}
		}
		if _, ok := c[e.Show][e.Season]; !ok {
			c[e.Show][e.Season] = map[int]string{}
		}
		c[e.Show][e.Season][e.Episode] = e.PlaylistItemID
	}
	l.episodeCache = c
}

func (p *Plex) EpisodeKey(show string, season, episode int) (string, error) {
	return "", errors.New("not yet implemented")
}

func (p *Plex) updateEpisodes(show string) error {
	if p.episodeIndex == nil {
		p.episodeIndex = map[string]map[int]map[int]string{}
	}
	if _, ok := p.episodeIndex[show]; !ok {
		p.episodeIndex[show] = map[int]map[int]string{}
	}
	return errors.New("not yet implemented")
}

// EpisodeKey returns the key of an episode, based on the title, season and episode
func (l *Playlist) EpisodeKey(title string, season, episode int) (string, error) {
	if l.episodeCache == nil {
		episodes, err := l.Episodes()
		if err != nil {
			return "", nil
		}
		l.updateEpisodeCache(&episodes)
	}
	if _, ok := l.episodeCache[title]; !ok {
		return "", fmt.Errorf("show not found in episode cache: %v", title)
	}
	if _, ok := l.episodeCache[title][season]; !ok {
		return "", fmt.Errorf("season not found in episode cache: %v: s%02d", title, season)
	}
	key, ok := l.episodeCache[title][season][episode]
	if !ok {
		return "", fmt.Errorf("episode not found in episode cache: %v s%02de%02d", title, season, episode)
	}
	return key, nil
}

// DeleteItem removes an item from the given playlist
func (l *Playlist) DeleteItem(keys ...string) error {
	for _, k := range keys {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%v/playlists/%v/items/%v", l.p.baseURL, l.ID, k), nil)
		if err != nil {
			return err
		}
		var got struct{}
		if err := l.p.sendRequest(req, &got); err != nil {
			return err
		}
	}
	return nil
}

// Episodes returns a new playlist by the name
func (l *Playlist) Episodes() (EpisodeList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/playlists/%v/items", l.p.baseURL, l.ID), nil)
	if err != nil {
		return nil, err
	}
	var plr PlaylistResponse
	if err := l.p.sendRequest(req, &plr); err != nil {
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

		var viewedInt int64
		if item.LastViewedAt != "" {
			var err error
			viewedInt, err = strconv.ParseInt(item.LastViewedAt, 10, 64)
			if err != nil {
				return nil, err
			}
		}
		viewed := time.Unix(viewedInt, 0)
		ret = append(ret, Episode{
			ID:             item.RatingKey,
			PlaylistItemID: item.PlaylistItemID,
			Title:          item.Title,
			Show:           item.GrandparentTitle,
			Season:         parentI,
			Episode:        index,
			Watched:        &viewed,
			p:              l.p,
		})
	}

	return ret, nil
}

func (p *Plex) sendRequest(req *http.Request, v interface{}) error {
	// req.Header.Set("Content-Type", "application/json; charset=utf-8")
	// req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("X-Plex-Token", p.token)

	if p.printCurl {
		command, _ := http2curl.GetCurlCommand(req)
		fmt.Fprintf(os.Stderr, "%v\n", command)
	}

	res, err := p.client.Do(req)
	if err != nil {
		return err
	}

	defer dclose(res.Body)

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	// fmt.Println(string(content))

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		if res.StatusCode == http.StatusTooManyRequests {
			return fmt.Errorf("too many requests.  Check rate limit and make sure the userAgent is set right")
		}
		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	if err = xml.NewDecoder(bytes.NewReader(content)).Decode(&v); err != nil {
		return err
	}

	return nil
}

type errorResponse struct {
	Status    string `json:"status"`
	ErrorType string `json:"errorType"`
	Error     string `json:"error"`
	Message   string `json:"message,omitempty"`
}

func dclose(c io.Closer) {
	if err := c.Close(); err != nil {
		slog.Error("error closing item", "error", err)
	}
}
