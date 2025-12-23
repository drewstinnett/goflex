package goflex

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"time"
)

// Episode represents an episode of television
type Episode struct {
	ID             int
	PlaylistItemID int
	Title          string
	Show           ShowTitle
	Season         SeasonNumber
	Episode        EpisodeNumber
	Watched        *time.Time
	ViewCount      int
	ViewOffset     *time.Duration
	Duration       time.Duration
}

// EpisodeNumber is the number of an episode
type EpisodeNumber int

// EpisodeMap is a map of episode numbers to episode objects
type EpisodeMap map[EpisodeNumber]*Episode

// List returns a list of seasons sorted by season number
func (e EpisodeMap) List() EpisodeList {
	seasonKeys := make([]int, len(e))
	i := 0
	for k := range e {
		seasonKeys[i] = int(k)
		i++
	}
	sort.Ints(seasonKeys)
	ret := make(EpisodeList, len(e))
	for idx, k := range seasonKeys {
		ret[idx] = *e[EpisodeNumber(k)]
	}
	return ret
}

// EpisodeFilter defines the filters the returned episodes
type EpisodeFilter struct {
	Show           ShowTitle    `yaml:"show"`
	EarliestSeason SeasonNumber `yaml:"earliest_season"`
	LatestSeason   SeasonNumber `yaml:"latest_season"`
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
func (l EpisodeList) Seasons(start, end SeasonNumber) EpisodeList {
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

/*
// ids returns a list of ids for the episodes
func (l EpisodeList) ids() []int {
	ret := make([]int, len(l))
	for idx, item := range l {
		ret[idx] = item.ID
	}
	return ret
}
*/

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
	switch e.Show {
	case "":
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
	return fmt.Sprintf("%v:%v:%v", e.Show, e.Season, e.Episode)
}

// Remaining returns the remaining time for the episode
func (e Episode) Remaining() time.Duration {
	if e.ViewOffset == nil {
		return e.Duration
	}
	return e.Duration - *e.ViewOffset
}

func episodeWithMetadata(m Metadata) (*Episode, error) {
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
		Show:      ShowTitle(m.GrandparentTitle),
		Season:    SeasonNumber(m.ParentIndex),
		ViewCount: m.ViewCount,
		Episode:   EpisodeNumber(m.Index),
	}, nil
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
