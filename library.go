package goflex

import (
	"fmt"
	"net/http"
	"strconv"
)

// LibraryService defines how to act against the Library service endpoints
type LibraryService interface {
	List() (map[string]*Library, error)
}

// LibraryServiceOp implements the LibraryService
type LibraryServiceOp struct {
	p     *Plex
	cache map[string]*Library
}

// List returns a list of libraries on the server
func (svc *LibraryServiceOp) List() (map[string]*Library, error) {
	if svc.cache == nil {
		if err := svc.updateLibraryCache(); err != nil {
			return nil, err
		}
	}
	return svc.cache, nil
}

// Library is a collection of movies or shows in Plex
type Library struct {
	ID        int
	Title     string
	Type      LibraryType
	showCache map[string]*Show
	p         *Plex
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
	}
	panic("unknown library type: " + s)
}

func (svc *LibraryServiceOp) updateLibraryCache() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/library/sections/", svc.p.baseURL), nil)
	if err != nil {
		return err
	}
	var lr LibraryResponse
	if err := svc.p.sendRequest(req, &lr); err != nil {
		return err
	}
	svc.cache = map[string]*Library{}
	for _, libd := range lr.Directory {
		id, err := strconv.Atoi(libd.Key)
		if err != nil {
			return err
		}
		svc.cache[libd.Title] = &Library{
			ID:    id,
			Title: libd.Title,
			Type:  stringToLibraryType(libd.Type),
			p:     svc.p,
		}
	}
	return nil
}
