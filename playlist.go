package plexrando

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
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
	List() (map[string]*Playlist, error)
	GetWithName(string) (*Playlist, error)
	Create(string, PlaylistType, bool) error
	GetOrCreate(string, PlaylistType, bool) (*Playlist, bool, error)
	Exists(string) bool
}

// PlaylistServiceOp is the operator for the PlaylistService
type PlaylistServiceOp struct {
	cache map[string]*Playlist
	p     *Plex
}

// GetOrCreate returns a playlist with the given title, or creates a new one with given kind and smart options
func (p *PlaylistServiceOp) GetOrCreate(title string, kind PlaylistType, smart bool) (*Playlist, bool, error) {
	if p.cache == nil {
		if err := p.updateCache(); err != nil {
			return nil, false, err
		}
	}
	var created bool
	if !p.Exists(title) {
		if err := p.Create(title, kind, smart); err != nil {
			return nil, false, err
		}
		created = true
	}
	got, err := p.GetWithName(title)
	if err != nil {
		return nil, created, err
	}
	return got, created, nil
}

// Exists returns true if a playlist already exists
func (p *PlaylistServiceOp) Exists(n string) bool {
	if p.cache == nil {
		if err := p.updateCache(); err != nil {
			panic(err)
		}
	}
	_, ok := p.cache[n]
	return ok
}

// Create creates a new playlist
func (p *PlaylistServiceOp) Create(title string, kind PlaylistType, smart bool) error {
	if p.Exists(title) {
		return errors.New("playlist with title already exists: " + title)
	}
	smartInt := 0
	if smart {
		smartInt = 1
	}
	req, err := http.NewRequest("POST",
		fmt.Sprintf(
			"%v/playlists?type=%v&title=%v&smart=%v&uri=server://%v/com.plexapp.plugins.library/", p.p.baseURL, kind, url.QueryEscape(title), smartInt, p.p.serverID),
		nil)
	if err != nil {
		return err
	}
	var res CreatePlaylistResponse
	if err := p.p.sendRequest(req, &res); err != nil {
		return err
	}
	return p.updateCache()
}

// List lists out playlists on a plex server
func (p *PlaylistServiceOp) List() (map[string]*Playlist, error) {
	if p.cache == nil {
		if err := p.updateCache(); err != nil {
			return nil, err
		}
	}
	return p.cache, nil
}

// GetWithName returns a playlist by name
func (p *PlaylistServiceOp) GetWithName(n string) (*Playlist, error) {
	if p.cache == nil {
		if err := p.updateCache(); err != nil {
			return nil, err
		}
	}
	got, ok := p.cache[n]
	if !ok {
		return nil, errors.New("playlist not found with name: " + n)
	}
	return got, nil
}

func (p *PlaylistServiceOp) updateCache() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/playlists", p.p.baseURL), nil)
	if err != nil {
		return err
	}
	var pr PlaylistsResponse
	if err := p.p.sendRequest(req, &pr); err != nil {
		return err
	}

	p.cache = map[string]*Playlist{}
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
		p.cache[item.Title] = &Playlist{
			ID:       id,
			Title:    item.Title,
			GUID:     item.GUID,
			Duration: time.Duration(second) * time.Second,
			p:        p.p,
		}
	}
	return nil
}

// Playlist returns a new playlist by the name
func (p *Plex) Playlist(s string) (*Playlist, error) {
	pl, ok := p.playlistMap[s]
	if !ok {
		return nil, fmt.Errorf("playlist not found. requested: %v, available: %v", s, p.playlistMap)
	}
	return &pl, nil
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
		viewed := time.Unix(viewedInt, 0)
		ret = append(ret, Episode{
			ID:             id,
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

// EpisodesDeprecated returns a new playlist by the name
func (l *Playlist) EpisodesDeprecated() (EpisodeList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/playlists/%v/items", l.p.baseURL, l.IDDeprecated), nil)
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
			DeprecatedID:   item.RatingKey,
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

// DeleteItem removes an item from the given playlist
func (l *Playlist) DeleteItem(keys ...string) error {
	for _, k := range keys {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%v/playlists/%v/items/%v", l.p.baseURL, l.IDDeprecated, k), nil)
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

// EpisodeKey returns the key of an episode, based on the title, season and episode
func (l *Playlist) EpisodeKey(title string, season, episode int) (string, error) {
	if l.episodeCache == nil {
		episodes, err := l.EpisodesDeprecated()
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
