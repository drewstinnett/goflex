package goflex

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"time"
)

// SessionService describes how to interact with sessions
type SessionService interface {
	All(*time.Time, ...ShowTitle) (EpisodeList, error)
	ActiveEpisodes(...ShowTitle) (EpisodeList, error)
	HistoryEpisodes(*time.Time, ...ShowTitle) (EpisodeList, error)
}

// SessionServiceOp is the operator for the session service
type SessionServiceOp struct {
	p *Plex
}

// All returns active and history sessions
func (s SessionServiceOp) All(since *time.Time, shows ...ShowTitle) (EpisodeList, error) {
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

func (svc *SessionServiceOp) historyEpisodes() (EpisodeList, error) {
	var res HistorySessionResponse
	if err := svc.p.sendRequestXML(mustNewRequest("GET",
		fmt.Sprintf("%v/status/sessions/history/all", svc.p.baseURL)),
		&res,
		&cacheConfig{
			prefix: "history-episodes",
			ttl:    time.Minute * 10,
		}); err != nil {
		return nil, err
	}
	ret := EpisodeList{}
	for _, item := range res.Video {
		if (item.Type != MediaTypeEpisode) || (item.RatingKey == "") {
			continue
		}

		viewedAt, err := dateFromUnixString(item.ViewedAt)
		if err != nil {
			return nil, err
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
			Show:    ShowTitle(item.GrandparentTitle),
			Season:  SeasonNumber(season),
			Episode: EpisodeNumber(index),
			Watched: viewedAt,
		})
	}
	return ret, nil
}

// HistoryEpisodes returns all episodes in the history. Given a list of shows, only returns watched episodes of those shows.
// Filter based on shows. Pass in a nil time.Time to return all times
func (svc *SessionServiceOp) HistoryEpisodes(since *time.Time, shows ...ShowTitle) (EpisodeList, error) {
	items, err := svc.historyEpisodes()
	if err != nil {
		return nil, err
	}
	if since == nil {
		return nil, errors.New("since cannot be nil")
	}
	ret := EpisodeList{}
	for _, item := range items {
		// Skip if we are listing shows, and this is not one of those shows
		switch {
		case (len(shows) > 0) && !slices.Contains(shows, item.Show):
			continue
		case item.Watched == nil:
			continue
		case item.Watched.Before(*since):
			slog.Debug("skipping because of since", "item", item, "since", *since)
			continue
		default:
			slog.Debug("adding item to ret", "item", item)
			ret = append(ret, item)
		}
	}

	return ret, nil
}

// ActiveEpisodes returns the active episodes in the session
func (s SessionServiceOp) ActiveEpisodes(shows ...ShowTitle) (EpisodeList, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/status/sessions", s.p.baseURL), nil)
	if err != nil {
		return nil, err
	}
	var res ActiveSessionsResponse
	if err := s.p.sendRequestXML(req, &res, &cacheConfig{prefix: "active-episodes", ttl: time.Minute * 5}); err != nil {
		return nil, err
	}
	ret := EpisodeList{}
	for _, item := range res.Video {
		if item.Type != MediaTypeEpisode {
			continue
		}
		if (len(shows) > 0) && !slices.Contains(shows, ShowTitle(item.GrandparentTitle)) {
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
			Show:    ShowTitle(item.GrandparentTitle),
			Season:  SeasonNumber(season),
			Episode: EpisodeNumber(index),
			// p:       s.p,
		})
	}

	return ret, nil
}
