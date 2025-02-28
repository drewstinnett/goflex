package goflex

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"time"
)

// SessionService describes how to interact with sessions
type SessionService interface {
	All(*time.Time, ...string) (EpisodeList, error)
	ActiveEpisodes(...string) (EpisodeList, error)
	HistoryEpisodes(*time.Time, ...string) (EpisodeList, error)
}

// SessionServiceOp is the operator for the session service
type SessionServiceOp struct {
	p *Plex
}

// All returns active and history sessions
func (s SessionServiceOp) All(since *time.Time, shows ...string) (EpisodeList, error) {
	ret, err := s.ActiveEpisodes(shows...)
	if err != nil {
		return nil, err
	}
	history, err := s.HistoryEpisodes(since, shows...)
	if err != nil {
		return nil, err
	}
	ret = append(ret, history...)
	return ret, nil
}

// HistoryEpisodes returns all episodes in the history. Given a list of shows, only returns watched episodes of those shows.
// Filter based on shows. Pass in a nil time.Time to return all times
func (s SessionServiceOp) HistoryEpisodes(since *time.Time, shows ...string) (EpisodeList, error) {
	var res HistorySessionResponse
	if err := s.p.sendRequestXML(mustNewRequest("GET", fmt.Sprintf("%v/status/sessions/history/all", s.p.baseURL)), &res); err != nil {
		return nil, err
	}
	ret := EpisodeList{}
	for _, item := range res.Video {
		if ((len(shows) > 0) && !slices.Contains(shows, item.GrandparentTitle)) ||
			(item.Type != MediaTypeEpisode) ||
			(item.RatingKey == "") {
			continue
		}

		viewedAt, err := dateFromUnixString(item.ViewedAt)
		if err != nil {
			return nil, err
		}
		if !viewedAt.IsZero() && (since != nil) && viewedAt.Before(*since) {
			continue
		}
		id, err := strconv.Atoi(item.RatingKey)
		if err != nil {
			return nil, err
		}
		index, err := strconv.Atoi(item.Index)
		if err != nil {
			return nil, err
		}
		season, err := strconv.Atoi(item.ParentIndex)
		if err != nil {
			return nil, err
		}

		// Skip items we have already added
		if slices.Contains(ret.ids(), id) {
			continue
		}

		ret = append(ret, Episode{
			ID:      id,
			Title:   item.Title,
			Show:    item.GrandparentTitle,
			Season:  season,
			Episode: index,
			Watched: viewedAt,
			p:       s.p,
		})
	}

	return ret, nil
}

// ActiveEpisodes returns the active episodes in the session
func (s SessionServiceOp) ActiveEpisodes(shows ...string) (EpisodeList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/status/sessions", s.p.baseURL), nil)
	if err != nil {
		return nil, err
	}
	var res ActiveSessionsResponse
	if err := s.p.sendRequestXML(req, &res); err != nil {
		return nil, err
	}
	ret := EpisodeList{}
	for _, item := range res.Video {
		if item.Type != MediaTypeEpisode {
			continue
		}
		if (len(shows) > 0) && !slices.Contains(shows, item.GrandparentTitle) {
			continue
		}
		id, err := strconv.Atoi(item.RatingKey)
		if err != nil {
			return nil, err
		}
		index, err := strconv.Atoi(item.Index)
		if err != nil {
			return nil, err
		}
		season, err := strconv.Atoi(item.ParentIndex)
		if err != nil {
			return nil, err
		}
		ret = append(ret, Episode{
			ID:      id,
			Title:   item.Title,
			Show:    item.GrandparentTitle,
			Season:  season,
			Episode: index,
			p:       s.p,
		})
	}

	return ret, nil
}
