package goflex

import (
	"slices"
	"strconv"
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
	Exists(ShowTitle) (bool, error)
	Match(ShowTitle) (ShowList, error)
	StrictMatch(ShowTitle) (ShowList, error)
	Seasons(Show) (*SeasonMap, error)
	SeasonsSorted(Show) (SeasonList, error)
	EpisodesWithFilter(ShowList, EpisodeFilter) (EpisodeList, error)
	Episodes(ShowList) (EpisodeList, error)
	SeasonEpisodes(*Season) (EpisodeMap, error)
}

// ShowMap maps show titles to shows
type ShowMap map[ShowTitle]*Show

// ShowTitle is the title of a show
type ShowTitle string

// Show represents a TV show in plex
type Show struct {
	ID    int
	Title ShowTitle
}

// Season represents a season in a TV show
type Season struct {
	ID    int
	Index int
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
