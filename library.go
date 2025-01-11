package plexrando

import (
	"fmt"
	"net/http"
	"strconv"
)

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

// Libraries returns a list of libraries on the server
func (p *Plex) Libraries() (map[string]Library, error) {
	if p.libraryCache == nil {
		if err := p.updateLibraryCache(); err != nil {
			return nil, err
		}
	}
	return p.libraryCache, nil
}

func (p *Plex) updateLibraryCache() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/library/sections/", p.baseURL), nil)
	if err != nil {
		return err
	}
	var lr LibraryResponse
	if err := p.sendRequest(req, &lr); err != nil {
		return err
	}
	p.libraryCache = map[string]Library{}
	for _, libd := range lr.Directory {
		id, err := strconv.Atoi(libd.Key)
		if err != nil {
			return err
		}
		p.libraryCache[libd.Title] = Library{
			ID:    id,
			Title: libd.Title,
			Type:  stringToLibraryType(libd.Type),
			p:     p,
		}
	}
	return nil
}
