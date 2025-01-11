package plexrando

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

// Show represents a TV show in plex
type Show struct {
	ID          int
	Title       string
	seasonCache map[int]*Season
	p           *Plex
}

// Season represents a season in a TV show
type Season struct {
	ID           int
	Index        int
	episodeCache map[int]*Episode
	p            *Plex
}

// MatchShows returns a list of shows that match a given name exactly
func (p *Plex) MatchShows(name string) ([]*Show, error) {
	libs, err := p.Libraries()
	if err != nil {
		return nil, err
	}
	ret := []*Show{}
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
			slog.Debug("skipping", "reason", "no rating key", "item", item)
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
			slog.Debug("skipping", "reason", "no rating key", "item", item)
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
