package goflex

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// LibraryService defines how to act against the Library service endpoints
type LibraryService interface {
	List() (LibraryMap, error)
	Shows(Library) (ShowMap, error)
}

// LibraryServiceOp implements the LibraryService
type LibraryServiceOp struct {
	p *Flex
}

// List returns a list of libraries on the server
func (svc *LibraryServiceOp) List() (LibraryMap, error) {
	var lr LibraryResponse
	if err := svc.p.sendRequestXML(mustNewRequest(http.MethodGet, fmt.Sprintf("%v/library/sections/", svc.p.baseURL)), &lr, &cacheConfig{prefix: "library-list", ttl: time.Minute * 60}); err != nil {
		return nil, err
	}
	ret := LibraryMap{}
	for _, libd := range lr.Directory {
		id, err := strconv.Atoi(libd.Key)
		if err != nil {
			return nil, err
		}
		ret[LibraryTitle(libd.Title)] = &Library{
			ID:    id,
			Title: libd.Title,
			Type:  stringToLibraryType(libd.Type),
		}
	}
	return ret, nil
}

// LibraryTitle just represents the title of the Library
type LibraryTitle string

// LibraryMap is a map of all libraries
type LibraryMap map[LibraryTitle]*Library

// Library is a collection of movies or shows in Plex
type Library struct {
	ID    int
	Title string
	Type  LibraryType
}

// LibraryType is the type of library
type LibraryType string

var (
	// ShowType represents a TV show library
	ShowType LibraryType = "show"
	// ArtistType represents a music library
	ArtistType LibraryType = "artist"
	// MovieType represents a movie library
	MovieType LibraryType = "movie"
)

// SearchType is the type of library search
type SearchType int

var (
	// SearchTypeMovie searches for movies
	SearchTypeMovie SearchType = 1
	// SearchTypeShow searches for TV shows
	SearchTypeShow SearchType = 2
	// SearchTypeEpisode searches for episodes
	SearchTypeEpisode SearchType = 4
	// SearchTypeArtist searches for music
	SearchTypeArtist SearchType = 8
)

func stringToLibraryType(s string) LibraryType {
	switch s {
	case "show":
		return ShowType
	case "artist":
		return ArtistType
	case "movie":
		return MovieType
	default:
		panic("unknown library type: " + s)
	}
}

// Shows returns all shows in a given library
func (svc *LibraryServiceOp) Shows(l Library) (ShowMap, error) {
	if l.Type != ShowType {
		return nil, errors.New("library is not a show library")
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%v/library/sections/%v/all", svc.p.baseURL, l.ID), nil)
	if err != nil {
		return nil, err
	}
	var sr ShowsResponse
	if err := svc.p.sendRequestXML(req, &sr, &cacheConfig{prefix: "shows", ttl: time.Minute * 5}); err != nil {
		return nil, err
	}
	ret := ShowMap{}
	for _, item := range sr.Directory {
		id, err := strconv.Atoi(item.RatingKey)
		if err != nil {
			return nil, err
		}
		ret[ShowTitle(item.Title)] = &Show{
			ID:    id,
			Title: ShowTitle(item.Title),
		}
	}
	return ret, nil
}
