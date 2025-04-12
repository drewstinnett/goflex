package goflex

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"time"
)

const (
	// MediaTypeEpisode is the string for "episode"
	MediaTypeEpisode string = "episode"
	// MediaTypeShow represents a show
	MediaTypeShow string = "show"
)

// ShowService describes how the show api behaves. This is a meta service where
// I'm putting some business logic on top of the API for tv shows.
type ShowService interface {
	Exists(ShowTitle) (bool, error)
	Match(ShowTitle) (ShowList, error)
	StrictMatch(ShowTitle) (ShowList, error)
	Seasons(Show) (*SeasonMap, error)
	SeasonsSorted(Show) (SeasonList, error)
	EpisodesWithFilter(ShowList, EpisodeFilter) (EpisodeList, error)
	Episodes(ShowList) (EpisodeList, error)
	Episode(ShowTitle, SeasonNumber, EpisodeNumber) (*Episode, error)
	SeasonEpisodes(*Season) (EpisodeMap, error)
}

// ShowServiceOp implements the ShowService operator.
type ShowServiceOp struct {
	p               *Flex
	cacheDeprecated ShowList
	// seasonCacheDeprecated map[ShowTitle]*SeasonMap
}

// Episode returns the episode for a given show, season, and episode number.
func (svc *ShowServiceOp) Episode(title ShowTitle, seasonNo SeasonNumber, episodeNo EpisodeNumber) (*Episode, error) {
	matches, err := svc.Match(title)
	if err != nil {
		return nil, err
	}
	for _, match := range matches {
		seasons, err := svc.Seasons(*match)
		if err != nil {
			return nil, err
		}
		for _, season := range *seasons {
			if season.Index != seasonNo {
				continue
			}
			episodes, err := svc.SeasonEpisodes(season)
			if err != nil {
				return nil, err
			}
			for _, episode := range episodes {
				if episode.Episode == episodeNo {
					return episode, nil
				}
			}
		}
	}
	return nil, errors.New("episode not found")
}

// Seasons returns the seasons for a given show.
func (svc *ShowServiceOp) Seasons(show Show) (*SeasonMap, error) {
	if show.Title == "" {
		return nil, errors.New("show.Title must not be empty")
	}
	var sr seasonsResponse
	if err := svc.p.sendRequestJSON(mustNewRequest(http.MethodGet, fmt.Sprintf("%v/library/metadata/%v/children", svc.p.baseURL, show.ID)), &sr, &cacheConfig{prefix: "seasons-" + string(show.Title), ttl: time.Hour * 1}); err != nil {
		return nil, fmt.Errorf("error sending json request: %w", err)
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
			Index: SeasonNumber(item.Index),
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
		if name == item.Title {
			return true, nil
		}
	}
	return false, nil
}

// StrictMatch returns an error if no shows are matched.
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

// Match returns shows with the given name.
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
		if show.Title == name {
			ret = append(ret, show)
		}
	}
	return ret, nil
}

// EpisodesWithFilter filters a shows episodes based on the given filter.
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

// Episodes returns episodes in a show list.
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

// SeasonEpisodes returns a list of episodes for a given season.
func (svc *ShowServiceOp) SeasonEpisodes(s *Season) (EpisodeMap, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%v/library/metadata/%v/children", svc.p.baseURL, s.ID),
		nil,
	)
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

		// Figure out how long the episode is.
		index, err := strconv.Atoi(item.Index)
		if err != nil {
			return nil, err
		}
		e, err := episodeWithVideo(item)
		if err != nil {
			return nil, err
		}
		ret[EpisodeNumber(index)] = e
	}
	return ret, nil
}

// SeasonsSorted returns a list of seasons sorted by season number.
func (svc *ShowServiceOp) SeasonsSorted(s Show) (SeasonList, error) {
	m, err := svc.Seasons(s)
	if err != nil {
		return nil, err
	}
	return m.sorted(), nil
}

// ShowMap maps show titles to shows.
type ShowMap map[ShowTitle]*Show

// ShowTitle is the title of a show.
type ShowTitle string

// Show represents a TV show in plex.
type Show struct {
	ID    int
	Title ShowTitle
}

// Season represents a season in a TV show.
type Season struct {
	ID    int
	Index SeasonNumber
}

// SeasonList is a list of Seasons
type SeasonList []*Season

// SeasonNumber represents the number in a given season. 0 is for Specials
type SeasonNumber int

// SeasonMap is a map of season numbers to season objects
type SeasonMap map[SeasonNumber]*Season

// sorted returns a list of seasons sorted by season number
func (s SeasonMap) sorted() SeasonList {
	seasonKeys := make([]SeasonNumber, len(s))
	i := 0
	for k := range s {
		seasonKeys[i] = k
		i++
	}
	slices.Sort(seasonKeys)
	ret := make(SeasonList, len(s))
	for idx, k := range seasonKeys {
		ret[idx] = s[k]
	}
	return ret
}

// ShowList represents multiple shows
type ShowList []*Show

func showWith(m Metadata) (*Show, error) {
	id, err := strconv.Atoi(m.RatingKey)
	if err != nil {
		return nil, err
	}
	return &Show{
		ID:    id,
		Title: ShowTitle(m.Title),
	}, nil
}

type seasonsResponse struct {
	Mediacontainer struct {
		Size                int    `json:"size"`
		Allowsync           bool   `json:"allowSync"`
		Art                 string `json:"art"`
		Identifier          string `json:"identifier"`
		Key                 string `json:"key"`
		Librarysectionid    int    `json:"librarySectionID"`
		Librarysectiontitle string `json:"librarySectionTitle"`
		Librarysectionuuid  string `json:"librarySectionUUID"`
		Mediatagprefix      string `json:"mediaTagPrefix"`
		Mediatagversion     int    `json:"mediaTagVersion"`
		Nocache             bool   `json:"nocache"`
		Parentindex         int    `json:"parentIndex"`
		Parenttitle         string `json:"parentTitle"`
		Parentyear          int    `json:"parentYear"`
		Summary             string `json:"summary"`
		Theme               string `json:"theme"`
		Thumb               string `json:"thumb"`
		Title1              string `json:"title1"`
		Title2              string `json:"title2"`
		Viewgroup           string `json:"viewGroup"`
		Directory           []struct {
			Leafcount       int    `json:"leafCount"`
			Thumb           string `json:"thumb"`
			Viewedleafcount int    `json:"viewedLeafCount"`
			Key             string `json:"key"`
			Title           string `json:"title"`
		} `json:"Directory"`
		Metadata []struct {
			RatingKey       string `json:"ratingKey"`
			Key             string `json:"key"`
			Parentratingkey string `json:"parentRatingKey"`
			GUID            string `json:"guid"`
			Parentguid      string `json:"parentGuid"`
			Parentslug      string `json:"parentSlug"`
			Parentstudio    string `json:"parentStudio"`
			Type            string `json:"type"`
			Title           string `json:"title"`
			Parentkey       string `json:"parentKey"`
			Parenttitle     string `json:"parentTitle"`
			Summary         string `json:"summary"`
			Index           int    `json:"index"`
			Parentindex     int    `json:"parentIndex"`
			Viewcount       int    `json:"viewCount"`
			Skipcount       int    `json:"skipCount,omitempty"`
			Lastviewedat    int    `json:"lastViewedAt"`
			Parentyear      int    `json:"parentYear"`
			Thumb           string `json:"thumb"`
			Art             string `json:"art"`
			Parentthumb     string `json:"parentThumb"`
			Parenttheme     string `json:"parentTheme"`
			Leafcount       int    `json:"leafCount"`
			Viewedleafcount int    `json:"viewedLeafCount"`
			Addedat         int    `json:"addedAt"`
			Updatedat       int    `json:"updatedAt"`
			Image           []struct {
				Alt  string `json:"alt"`
				Type string `json:"type"`
				URL  string `json:"url"`
			} `json:"Image"`
			Ultrablurcolors struct {
				Topleft     string `json:"topLeft"`
				Topright    string `json:"topRight"`
				Bottomright string `json:"bottomRight"`
				Bottomleft  string `json:"bottomLeft"`
			} `json:"UltraBlurColors,omitempty"`
		} `json:"Metadata"`
	} `json:"MediaContainer"`
}
