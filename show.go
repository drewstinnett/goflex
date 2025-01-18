package plexrando

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"time"
)

// Show represents a TV show in plex
type Show struct {
	ID          int
	Title       string
	seasonCache map[int]*Season
	p           *Plex
}

// Episode represents an episode of television
type Episode struct {
	ID             int
	PlaylistItemID int
	Title          string
	Show           string
	Season         int
	Episode        int
	Watched        *time.Time
	ViewCount      int
	p              *Plex
}

// Season represents a season in a TV show
type Season struct {
	ID           int
	Index        int
	episodeCache map[int]*Episode
	p            *Plex
}

// ShowList represents multiple shows
type ShowList []*Show

// Episodes returns episodes in a show list
func (s *ShowList) Episodes() (EpisodeList, error) {
	ret := EpisodeList{}
	for _, show := range *s {
		seasons, err := show.Seasons()
		if err != nil {
			return nil, err
		}
		for _, season := range seasons {
			episodes, err := season.Episodes()
			if err != nil {
				return nil, err
			}
			for _, episode := range episodes {
				ret = append(ret, *episode)
			}
		}
	}
	return ret, nil
}

// MatchShows returns a list of shows that match a given name exactly
func (p *Plex) MatchShows(name string) (ShowList, error) {
	libs, err := p.Libraries()
	if err != nil {
		return nil, err
	}
	ret := ShowList{}
	for _, lib := range libs {
		if lib.Type != ShowType {
			continue
		}
		shows, err := lib.Shows()
		if err != nil {
			return nil, err
		}
		for _, show := range shows {
			if show.Title == name {
				ret = append(ret, show)
			}
		}
	}
	return ret, nil
}

// Episodes returns a list of episodes for a given season
func (s *Season) Episodes() (map[int]*Episode, error) {
	if s.episodeCache == nil {
		if err := s.updateEpisodeCache(); err != nil {
			return nil, err
		}
	}
	return s.episodeCache, nil
}

func (s *Season) updateEpisodeCache() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/library/metadata/%v/children", s.p.baseURL, s.ID), nil)
	if err != nil {
		return err
	}
	var er EpisodesResponse
	if err := s.p.sendRequest(req, &er); err != nil {
		return err
	}
	s.episodeCache = map[int]*Episode{}
	for _, item := range er.Video {
		if item.RatingKey == "" {
			continue
		}
		id, err := strconv.Atoi(item.RatingKey)
		if err != nil {
			return err
		}
		index, err := strconv.Atoi(item.Index)
		if err != nil {
			return err
		}
		s.episodeCache[index] = &Episode{
			ID:      id,
			Title:   item.Title,
			Show:    item.GrandparentTitle,
			Season:  s.Index,
			Episode: index,
			p:       s.p,
		}
	}
	return nil
}

func (s *Show) updateSeasonCache() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/library/metadata/%v/children", s.p.baseURL, s.ID), nil)
	if err != nil {
		return err
	}
	var sr SeasonsResponse
	if err := s.p.sendRequest(req, &sr); err != nil {
		return err
	}
	s.seasonCache = map[int]*Season{}
	for _, item := range sr.Directory {
		if item.RatingKey == "" {
			continue
		}
		id, err := strconv.Atoi(item.RatingKey)
		if err != nil {
			return err
		}
		index, err := strconv.Atoi(item.Index)
		if err != nil {
			return err
		}
		s.seasonCache[index] = &Season{
			ID:    id,
			Index: index,
			p:     s.p,
		}
	}
	return nil
}

// Seasons returns a list of seasons in a given show
func (s *Show) Seasons() (map[int]*Season, error) {
	if s.seasonCache == nil {
		if err := s.updateSeasonCache(); err != nil {
			return nil, err
		}
	}
	return s.seasonCache, nil
}

// Shows returns a list of TV shows in a given library
func (l *Library) Shows() (map[string]*Show, error) {
	if l.Type != ShowType {
		return nil, errors.New("library is not a show library")
	}
	if l.showCache == nil {
		if err := l.updateShowCache(); err != nil {
			return nil, err
		}
	}
	return l.showCache, nil
}

func (l *Library) updateShowCache() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/library/sections/%v/all", l.p.baseURL, l.ID), nil)
	if err != nil {
		return err
	}
	var sr ShowsResponse
	if err := l.p.sendRequest(req, &sr); err != nil {
		return err
	}
	l.showCache = map[string]*Show{}
	for _, item := range sr.Directory {
		id, err := strconv.Atoi(item.RatingKey)
		if err != nil {
			return err
		}
		l.showCache[item.Title] = &Show{
			ID:    id,
			Title: item.Title,
			p:     l.p,
		}
	}
	return nil
}

func (p *Plex) episodeID(show string, season, episode int) (int, error) {
	shows, err := p.MatchShows(show)
	if err != nil {
		return 0, err
	}
	episodes, err := shows.Episodes()
	if err != nil {
		return 0, err
	}
	for _, item := range episodes {
		if (item.Season == season) && (item.Episode == episode) {
			return item.ID, nil
		}
	}
	return 0, errors.New("episode key not found")
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
	ids := s.ids()
	for _, item := range *l {
		if !slices.Contains(ids, item.ID) {
			r = append(r, item)
		} else {
			removed = append(removed, item)
		}
	}
	return r, removed
}

// ids returns a list of ids for the episodes
func (l EpisodeList) ids() []int {
	ret := make([]int, len(l))
	for idx, item := range l {
		ret[idx] = item.ID
	}
	return ret
}

// String fulfills the Stringer interface
func (e Episode) String() string {
	var ret string
	switch {
	case e.Show == "":
		ret = fmt.Sprintf("%v - %v", e.ID, e.Title)
	default:
		ret = fmt.Sprintf("%v - S%02dE%02d - %v", e.Show, e.Season, e.Episode, e.Title)
	}
	if e.Watched != nil {
		ret = fmt.Sprintf("%v (Watched: %v)", ret, e.Watched.Format("2006-01-02 15:04"))
	}
	return ret
}
