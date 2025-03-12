package goflex

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// ShowServiceOp implements the ShowService operator
type ShowServiceOp struct {
	p                     *Plex
	cacheDeprecated       ShowList
	seasonCacheDeprecated map[ShowTitle]*SeasonMap
}

// Seasons returns the seasons for a given show
func (svc *ShowServiceOp) Seasons(show Show) (*SeasonMap, error) {
	/*
		if svc.seasonCacheDeprecated == nil {
			svc.seasonCacheDeprecated = map[ShowTitle]*SeasonMap{}
		}
		if _, ok := svc.seasonCacheDeprecated[show.Title]; !ok {
			if err := svc.updateSeasonCache(show); err != nil {
				return nil, err
			}
		}
	*/
	var sr seasonsResponse
	if err := svc.p.sendRequestJSON(mustNewRequest("GET", fmt.Sprintf("%v/library/metadata/%v/children", svc.p.baseURL, show.ID)), &sr, toPTR(time.Hour*1)); err != nil {
		return nil, err
	}
	ret := SeasonMap{}
	for _, item := range sr.Mediacontainer.Metadata {
		if item.RatingKey == "" {
			continue
		}
		id, err := strconv.Atoi(item.RatingKey)
		if err != nil {
			return nil, err
		}
		ret[SeasonNumber(item.Index)] = &Season{
			ID:    id,
			Index: item.Index,
		}
	}
	return &ret, nil
}

func (svc *ShowServiceOp) updateCacheDeprecated() error {
	libs, err := svc.p.Library.List()
	if err != nil {
		return err
	}
	svc.cacheDeprecated = ShowList{}

	for _, lib := range libs {
		if lib.Type != ShowType {
			continue
		}
		shows, err := svc.p.Library.Shows(*lib)
		if err != nil {
			return err
		}
		for _, show := range shows {
			svc.cacheDeprecated = append(svc.cacheDeprecated, show)
		}
	}
	return nil
}

// Exists returns true if a show exists on the server
func (svc *ShowServiceOp) Exists(name ShowTitle) (bool, error) {
	if svc.cacheDeprecated == nil {
		if err := svc.updateCacheDeprecated(); err != nil {
			return false, err
		}
		/*
		 */
	}
	for _, item := range svc.cacheDeprecated {
		if ShowTitle(name) == item.Title {
			return true, nil
		}
	}
	return false, nil
}

// StrictMatch returns an error if no shows are matched
func (svc *ShowServiceOp) StrictMatch(name ShowTitle) (ShowList, error) {
	got, err := svc.Match(name)
	if err != nil {
		return nil, err
	}
	if len(got) == 0 {
		return nil, errors.New("no shows found matching: " + string(name))
	}
	return got, nil
}

// Match returns shows with the given name
func (svc *ShowServiceOp) Match(name ShowTitle) (ShowList, error) {
	if svc.cacheDeprecated == nil {
		if err := svc.updateCacheDeprecated(); err != nil {
			return nil, err
		}
		/*
			if err := svc.updateCache(); err != nil {
				return nil, err
			}
		*/
	}
	ret := ShowList{}
	for _, show := range svc.cacheDeprecated {
		if show.Title == ShowTitle(name) {
			ret = append(ret, show)
		}
	}
	return ret, nil
}

// EpisodesWithFilter filters a shows episodes based on the given filter
func (svc *ShowServiceOp) EpisodesWithFilter(s ShowList, f EpisodeFilter) (EpisodeList, error) {
	ret := EpisodeList{}
	for _, show := range s {
		// seasons, err := show.Seasons()
		seasons, err := svc.Seasons(*show)
		if err != nil {
			return nil, err
		}
		for _, season := range *seasons {
			if ((f.EarliestSeason != 0) && (season.Index < f.EarliestSeason)) ||
				((f.LatestSeason != 0) && (season.Index > f.LatestSeason)) {
				continue
			}
			// episodes, err := season.Episodes()
			episodes, err := svc.SeasonEpisodes(season)
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

// Episodes returns episodes in a show list
func (svc *ShowServiceOp) Episodes(s ShowList) (EpisodeList, error) {
	ret := EpisodeList{}
	for _, show := range s {
		// seasons, err := show.Seasons()
		seasons, err := svc.Seasons(*show)
		if err != nil {
			return nil, err
		}
		for _, season := range *seasons {
			episodes, err := svc.SeasonEpisodes(season)
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

// SeasonEpisodes returns a list of episodes for a given season
func (svc *ShowServiceOp) SeasonEpisodes(s *Season) (EpisodeMap, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/library/metadata/%v/children", svc.p.baseURL, s.ID), nil)
	if err != nil {
		return nil, err
	}
	var er EpisodesResponse
	if err := svc.p.sendRequestXML(req, &er, nil); err != nil {
		return nil, err
	}
	// s.episodeCache = map[EpisodeNumber]*Episode{}
	ret := EpisodeMap{}
	for _, item := range er.Video {
		if item.RatingKey == "" {
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
		vc := 0
		if item.ViewCount != "" {
			var err error
			if vc, err = strconv.Atoi(item.ViewCount); err != nil {
				return nil, err
			}
		}
		e := &Episode{
			ID:        id,
			Title:     item.Title,
			Show:      ShowTitle(item.GrandparentTitle),
			Season:    SeasonNumber(s.Index),
			ViewCount: vc,
			Episode:   EpisodeNumber(index),
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
		ret[EpisodeNumber(index)] = e
	}
	return ret, nil
}

// SeasonsSorted returns a list of seasons sorted by season number
func (svc *ShowServiceOp) SeasonsSorted(s Show) (SeasonList, error) {
	m, err := svc.Seasons(s)
	if err != nil {
		return nil, err
	}
	return m.sorted(), nil
}
