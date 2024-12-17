/*
Package plexrando does the randomization bits
*/
package plexrando

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"time"

	"github.com/LukeHagar/plexgo"
	"github.com/LukeHagar/plexgo/models/operations"
)

// Plex connects to our custom stuff
type Plex struct {
	serverID    string
	API         *plexgo.PlexAPI
	LibraryMap  map[string]int
	PlaylistMap map[string]Playlist
}

// Episode represents an episode of television
type Episode struct {
	ID      string // ID is pretty much just the RatingKey
	Title   string
	Show    string
	Season  int
	Episode int
	Watched *time.Time
	p       *Plex
}

// EpisodeList is multiple Episodes
type EpisodeList []Episode

// Runtime returns the total runtime of all episodes
func (l EpisodeList) Runtime() time.Duration {
	/*
		ret := 0
		for _, item := range l {
		}
	*/
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

// New uses functional options for a new plex
func New(opts ...func(*Plex)) *Plex {
	p := &Plex{}
	for _, opt := range opts {
		opt(p)
	}
	if err := p.init(); err != nil {
		panic(err)
	}
	return p
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
	p.PlaylistMap = map[string]Playlist{}
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
		p.PlaylistMap[*playlist.Title] = Playlist{
			ID:       k,
			URI:      *playlist.GUID,
			Duration: duration,
			p:        p,
		}
	}
	return nil
}

func (p *Plex) updateSections() error {
	p.LibraryMap = map[string]int{}
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
		p.LibraryMap[libd.Title] = id
	}
	return nil
}

// GetOrCreatePlaylist creates an empty playlist using the given name. If a playlist
// with that name already exists, it will be emptied
func (p *Plex) GetOrCreatePlaylist(s string) (*Playlist, bool, error) {
	var playlist Playlist
	var ok bool
	if playlist, ok = p.PlaylistMap[s]; !ok {
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
	lib := p.LibraryMap[library]
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

// PlaylistEpisodes returns a list of episodes for a show
func (p *Plex) PlaylistEpisodes(title string) (EpisodeList, error) {
	id, ok := p.PlaylistMap[title]
	if !ok {
		return nil, fmt.Errorf("unknown library: %v", title)
	}
	ret := EpisodeList{}
	res, err := p.API.Playlists.GetPlaylistContents(context.Background(), id.ID, operations.GetPlaylistContentsQueryParamTypeEpisode)
	if err != nil {
		return nil, err
	}
	for _, item := range res.Object.MediaContainer.Metadata {
		ret = append(ret, p.NewEpisode(item))
	}
	return ret, nil
}

// Episodes returns a list of episodes for a show
func (p *Plex) Episodes(library, show string) (EpisodeList, error) {
	id, ok := p.LibraryMap[library]
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

// Playlist is the important identifiers of a playlist
type Playlist struct {
	ID       float64
	URI      string
	Duration time.Duration
	p        *Plex
}

// AddEpisodes adds an episode or more to a given playlist
func (l Playlist) AddEpisodes(episodes EpisodeList) error {
	slog.Debug("Adding episodes back in", "count", len(episodes))
	for idx, episode := range episodes {
		i := float64(idx)
		if _, err := l.p.API.Playlists.AddPlaylistContents(context.Background(), l.ID, episode.URI(), &i); err != nil {
			return err
		}
		slog.Debug("Adding", "episode", episode.String())
	}
	return l.p.updatePlaylists()
}
