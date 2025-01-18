/*
Package plexrando does the randomization bits
*/
package plexrando

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"moul.io/http2curl/v2"
)

const version string = "0.1.0"

// Plex connects to our custom stuff
type Plex struct {
	baseURL   string
	token     string
	userAgent string
	printCurl bool
	client    *http.Client
	// API          *plexgo.PlexAPI
	initialize   bool
	libraryCache map[string]Library
	// libraryMap   map[string]int
	playlistMap map[string]Playlist
	// showLibraryMap map[string]string
	// episodeIndex map[string]map[int]map[int]string
	Playlists PlaylistService
	Sessions  SessionService
	Media     MediaService
	Server    ServerService
}

// New uses functional options for a new plex
func New(opts ...func(*Plex)) (*Plex, error) {
	p := &Plex{
		initialize: true,
		client:     http.DefaultClient,
		userAgent:  "goflex " + version,
	}
	for _, opt := range opts {
		opt(p)
	}
	if p.baseURL == "" {
		return nil, errors.New("must set plex baseurl")
	}
	if p.token == "" {
		return nil, errors.New("must set token")
	}

	// Set up a mainstream API attribute
	/*
		p.API = plexgo.New(
			plexgo.WithSecurity(p.token),
			plexgo.WithServerURL(p.baseURL),
			plexgo.WithClientID("313FF6D7-5795-45E3-874F-B8FCBFD5E587"),
			plexgo.WithClientName("plex-trueget"),
			plexgo.WithClientVersion(version),
		)
	*/
	p.Playlists = &PlaylistServiceOp{p: p}
	p.Sessions = &SessionServiceOp{p: p}
	p.Media = &MediaServiceOp{p: p}
	p.Server = &ServerServiceOp{p: p}
	if p.initialize {
		if err := p.init(); err != nil {
			panic(err)
		}
	}
	return p, nil
}

// WithoutInit skips the initialization step
func WithoutInit() func(*Plex) {
	return func(p *Plex) {
		p.initialize = false
	}
}

// WithHTTPClient sets the http client on a new Plex
func WithHTTPClient(c *http.Client) func(*Plex) {
	return func(p *Plex) {
		p.client = c
	}
}

// WithBaseURL sets the base url for a plex client
func WithBaseURL(s string) func(*Plex) {
	return func(p *Plex) {
		p.baseURL = s
	}
}

// WithToken sets the plex token for a client
func WithToken(s string) func(*Plex) {
	return func(p *Plex) {
		p.token = s
	}
}

// WithPrintCurl sets the curl debug printer to on
func WithPrintCurl() func(*Plex) {
	return func(p *Plex) {
		p.printCurl = true
	}
}

func (p *Plex) init() error {
	return nil
}

type playlistEpisodeCache map[string]map[int]map[int]int

// Playlist is the important identifiers of a playlist
type Playlist struct {
	IDDeprecated float64
	ID           int
	Title        string
	GUID         string
	Duration     time.Duration
	p            *Plex
	episodeCache playlistEpisodeCache
}

func (p *Plex) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("X-Plex-Token", p.token)
	req.Header.Set("User-Agent", p.userAgent)

	if p.printCurl {
		command, _ := http2curl.GetCurlCommand(req)
		fmt.Fprintf(os.Stderr, "%v\n", command)
	}

	res, err := p.client.Do(req)
	if err != nil {
		return err
	}

	defer dclose(res.Body)

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	// fmt.Println(string(content))

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		if res.StatusCode == http.StatusTooManyRequests {
			return fmt.Errorf("too many requests.  Check rate limit and make sure the userAgent is set right")
		}
		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	if err = xml.NewDecoder(bytes.NewReader(content)).Decode(&v); err != nil && err != io.EOF {
		return err
	}

	return nil
}

type errorResponse struct {
	Status    string `json:"status"`
	ErrorType string `json:"errorType"`
	Error     string `json:"error"`
	Message   string `json:"message,omitempty"`
}

func dclose(c io.Closer) {
	if err := c.Close(); err != nil {
		slog.Error("error closing item", "error", err)
	}
}

/*
func toPTR[V any](v V) *V {
	return &v
}
*/
