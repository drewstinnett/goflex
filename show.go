package goflex

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"sort"
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
// I'm putting some business logic on top of the API for tv shows
type ShowService interface {
	Exists(string) (bool, error)
	Match(string) (ShowList, error)
	StrictMatch(string) (ShowList, error)
	Seasons(Show) (*SeasonMap, error)
}

// ShowServiceOp implements the ShowService operator
type ShowServiceOp struct {
	p           *Plex
	cache       ShowList
	seasonCache map[string]*SeasonMap
}

// Seasons returns the seasons for a given show
func (svc *ShowServiceOp) Seasons(show Show) (*SeasonMap, error) {
	if svc.seasonCache == nil {
		svc.seasonCache = map[string]*SeasonMap{}
	}
	if _, ok := svc.seasonCache[show.Title]; !ok {
		if err := svc.updateSeasonCache(show); err != nil {
			return nil, err
		}
	}
	return svc.seasonCache[show.Title], nil
}

func (svc *ShowServiceOp) updateSeasonCache(s Show) error {
	var sr seasonsResponse
	if err := svc.p.sendRequestJSON(mustNewRequest("GET", fmt.Sprintf("%v/library/metadata/%v/children", svc.p.baseURL, s.ID)), &sr); err != nil {
		return err
	}
	svc.seasonCache = map[string]*SeasonMap{}
	svc.seasonCache[s.Title] = &SeasonMap{}
	sm := SeasonMap{}
	for _, item := range sr.Mediacontainer.Metadata {
		if item.RatingKey == "" {
			continue
		}
		id, err := strconv.Atoi(item.RatingKey)
		if err != nil {
			return err
		}
		sm[item.Index] = &Season{
			ID:    id,
			Index: item.Index,
			p:     s.p,
		}
	}
	svc.seasonCache[s.Title] = &sm
	return nil
}

/*
func (svc *ShowServiceOp) updateSeasonCacheDeprecated(s Show) error {
	var sr seasonsResponse
	if err := svc.p.sendRequestJSON(mustNewRequest("GET", fmt.Sprintf("%v/library/metadata/%v/children", svc.p.baseURL, s.ID)), &sr); err != nil {
		return err
	}
	s.seasonCache = map[int]*Season{}
	for _, item := range sr.Mediacontainer.Metadata {
		if item.RatingKey == "" {
			continue
		}
		id, err := strconv.Atoi(item.RatingKey)
		if err != nil {
			return err
		}
		s.seasonCache[item.Index] = &Season{
			ID:    id,
			Index: item.Index,
			p:     s.p,
		}
	}
	return nil
}
*/

func (svc *ShowServiceOp) updateCache() error {
	libs, err := svc.p.Library.List()
	if err != nil {
		return err
	}
	svc.cache = ShowList{}

	for _, lib := range libs {
		if lib.Type != ShowType {
			continue
		}
		shows, err := lib.Shows()
		if err != nil {
			return err
		}
		for _, show := range shows {
			svc.cache = append(svc.cache, show)
		}
	}
	return nil
}

// Exists returns true if a show exists on the server
func (svc *ShowServiceOp) Exists(name string) (bool, error) {
	if svc.cache == nil {
		if err := svc.updateCache(); err != nil {
			return false, err
		}
	}
	for _, item := range svc.cache {
		if name == item.Title {
			return true, nil
		}
	}
	return false, nil
}

// StrictMatch returns an error if no shows are matched
func (svc *ShowServiceOp) StrictMatch(name string) (ShowList, error) {
	got, err := svc.Match(name)
	if err != nil {
		return nil, err
	}
	if len(got) == 0 {
		return nil, errors.New("no shows found matching: " + name)
	}
	return got, nil
}

// Match returns shows with the given name
func (svc *ShowServiceOp) Match(name string) (ShowList, error) {
	if svc.cache == nil {
		if err := svc.updateCache(); err != nil {
			return nil, err
		}
	}
	ret := ShowList{}
	for _, show := range svc.cache {
		if show.Title == name {
			ret = append(ret, show)
		}
	}
	return ret, nil
}

// Show represents a TV show in plex
type Show struct {
	ID          int
	Title       string
	seasonCache SeasonMap
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
	episodeCache EpisodeMap
	p            *Plex
}

// SeasonList is a list of Seasons
type SeasonList []*Season

// EpisodeMap is a map of episode numbers to episode objects
type EpisodeMap map[int]*Episode

// List returns a list of seasons sorted by season number
func (e EpisodeMap) List() EpisodeList {
	seasonKeys := make([]int, len(e))
	i := 0
	for k := range e {
		seasonKeys[i] = k
		i++
	}
	sort.Ints(seasonKeys)
	ret := make(EpisodeList, len(e))
	for idx, k := range seasonKeys {
		ret[idx] = *e[k]
	}
	return ret
}

// SeasonMap is a map of season numbers to season objects
type SeasonMap map[int]*Season

// sorted returns a list of seasons sorted by season number
func (s SeasonMap) sorted() SeasonList {
	seasonKeys := make([]int, len(s))
	i := 0
	for k := range s {
		seasonKeys[i] = k
		i++
	}
	sort.Ints(seasonKeys)
	ret := make(SeasonList, len(s))
	for idx, k := range seasonKeys {
		ret[idx] = s[k]
	}
	return ret
}

// ShowList represents multiple shows
type ShowList []*Show

// EpisodeFilter defines the filters the returned episodes
type EpisodeFilter struct {
	Show           string
	EarliestSeason int
	LatestSeason   int
}

// EpisodesWithFilter filters a shows episodes based on the given filter
func (s *ShowList) EpisodesWithFilter(f EpisodeFilter) (EpisodeList, error) {
	ret := EpisodeList{}
	for _, show := range *s {
		seasons, err := show.Seasons()
		if err != nil {
			return nil, err
		}
		for _, season := range seasons {
			if ((f.EarliestSeason != 0) && (season.Index < f.EarliestSeason)) ||
				((f.LatestSeason != 0) && (season.Index > f.LatestSeason)) {
				continue
			}
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

// Episodes returns a list of episodes for a given season
func (s *Season) Episodes() (EpisodeMap, error) {
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
	if err := s.p.sendRequestXML(req, &er); err != nil {
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
		vc := 0
		if item.ViewCount != "" {
			var err error
			if vc, err = strconv.Atoi(item.ViewCount); err != nil {
				return err
			}
		}
		e := &Episode{
			ID:        id,
			Title:     item.Title,
			Show:      item.GrandparentTitle,
			Season:    s.Index,
			ViewCount: vc,
			Episode:   index,
			p:         s.p,
		}
		if item.LastViewedAt != "" {
			var viewedInt int64
			var err error
			viewedInt, err = strconv.ParseInt(item.LastViewedAt, 10, 64)
			if err != nil {
				return err
			}
			viewed := time.Unix(viewedInt, 0)
			e.Watched = &viewed
		}
		s.episodeCache[index] = e
	}
	return nil
}

func (s *Show) updateSeasonCache() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/library/metadata/%v/children", s.p.baseURL, s.ID), nil)
	if err != nil {
		return err
	}
	var sr SeasonsResponse
	if err := s.p.sendRequestXML(req, &sr); err != nil {
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
func (s *Show) Seasons() (SeasonMap, error) {
	if s.seasonCache == nil {
		if err := s.updateSeasonCache(); err != nil {
			return nil, err
		}
	}
	return s.seasonCache, nil
}

// SeasonsSorted returns a list of seasons sorted by season number
func (s *Show) SeasonsSorted() (SeasonList, error) {
	m, err := s.Seasons()
	if err != nil {
		return nil, err
	}
	return m.sorted(), nil
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
	if err := l.p.sendRequestXML(req, &sr); err != nil {
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
	shows, err := p.Shows.Match(show)
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

// EarliestWatched returns the episode watched the earliest in the list
func (l EpisodeList) EarliestWatched() *Episode {
	var ret *Episode
	for _, episode := range l {
		if episode.Watched != nil {
			if (ret == nil) || (ret.Watched.After(*episode.Watched)) {
				ret = &episode
			}
		}
	}
	return ret
}

// LatestWatched returns the episode watched the latest in the list
func (l EpisodeList) LatestWatched() *Episode {
	var ret *Episode
	for _, episode := range l {
		if episode.Watched != nil {
			if (ret == nil) || (!ret.Watched.After(*episode.Watched)) {
				ret = &episode
			}
		}
	}
	return ret
}

// WatchSpan returns the time between the earliest and latest watch
func (l EpisodeList) WatchSpan() (*time.Duration, error) {
	earliest := l.EarliestWatched()
	if earliest == nil {
		return nil, errors.New("could not find earliest watched")
	}

	latest := l.LatestWatched()
	if latest == nil {
		return nil, errors.New("could not find latest watched")
	}

	return toPTR(latest.Watched.Sub(fromPTR(earliest.Watched))), nil
}

// Runtime returns the total runtime of all episodes
func (l EpisodeList) Runtime() time.Duration {
	return 0
}

// Len returns the length of the list, to satisfy the sortable interface
func (l EpisodeList) Len() int {
	return len(l)
}

// Swap fulfills the sortable interface
func (l EpisodeList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

// Less determines if one episode appears sooner than another
func (l EpisodeList) Less(i, j int) bool {
	if l[i].Season != l[j].Season {
		return l[i].Season < l[j].Season
	}
	return l[i].Episode < l[j].Episode
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
	seenSlugs := []string{}
	r := EpisodeList{}
	removed := EpisodeList{}
	slugs := s.slugs()
	for _, item := range *l {
		slug := item.slug()
		if slices.Contains(seenSlugs, slug) {
			continue
		}
		if !slices.Contains(slugs, slug) {
			r = append(r, item)
		} else {
			removed = append(removed, item)
		}
		seenSlugs = append(seenSlugs, slug)
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

// slugs returns a list of slugs for the episodes
func (l EpisodeList) slugs() []string {
	ret := make([]string, len(l))
	for idx, item := range l {
		ret[idx] = item.slug()
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
	if e.ViewCount > 0 {
		ret = fmt.Sprintf("%v (ViewCount: %v)", ret, e.ViewCount)
	}
	return ret
}

// slug returns a unique string for a given episode
func (e Episode) slug() string {
	return fmt.Sprintf("%v:%v:%v", e.Season, e.Season, e.Episode)
}

func episodeWith(m Metadata) (*Episode, error) {
	if m.RatingKey == "" {
		return nil, errors.New("missing rating key")
	}
	id, err := strconv.Atoi(m.RatingKey)
	if err != nil {
		return nil, err
	}
	return &Episode{
		ID:        id,
		Title:     m.Title,
		Show:      m.GrandparentTitle,
		Season:    m.ParentIndex,
		ViewCount: m.ViewCount,
		Episode:   m.Index,
	}, nil
}

func showWith(m Metadata) (*Show, error) {
	id, err := strconv.Atoi(m.RatingKey)
	if err != nil {
		return nil, err
	}
	return &Show{
		ID:    id,
		Title: m.Title,
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
