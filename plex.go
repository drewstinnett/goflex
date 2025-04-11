/*
Package goflex does the randomization bits
*/
package goflex

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
	baseURL        string
	token          string
	userAgent      string
	printCurl      bool
	maxSleep       time.Duration
	client         *http.Client
	cache          cache
	Playlists      PlaylistService
	Sessions       SessionService
	Media          MediaService
	Server         ServerService
	Shows          ShowService
	Library        LibraryService
	Authentication AuthenticationService
}

// New uses functional options for a new plex
func New(opts ...func(*Plex)) (*Plex, error) {
	p := &Plex{
		client:    http.DefaultClient,
		userAgent: "goflex " + version,
		cache:     *NewCache(),
		maxSleep:  60 * time.Minute,
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

	p.Playlists = &PlaylistServiceOp{p: p}
	p.Sessions = &SessionServiceOp{p: p}
	p.Media = &MediaServiceOp{p: p}
	p.Server = &ServerServiceOp{p: p}
	p.Shows = &ShowServiceOp{p: p}
	p.Library = &LibraryServiceOp{p: p}
	p.Authentication = &AuthenticationServiceOp{p: p}
	return p, nil
}

// WithGCInterval sets the garbage collection interval
func WithGCInterval(i *time.Duration) func(*Plex) {
	return func(p *Plex) {
		p.cache = *NewCacheWithGC(fromPTR(i))
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

func (p *Plex) preprocessReq(req *http.Request) {
	req.Header.Set("X-Plex-Token", p.token)
	req.Header.Set("User-Agent", p.userAgent)

	if p.printCurl {
		command, _ := http2curl.GetCurlCommand(req)
		fmt.Fprintf(os.Stderr, "%v\n", command)
	}
}

const (
	jsonHeader string = "application/json"
	xmlHeader  string = "application/xml"
)

func (p *Plex) sendRequestXML(req *http.Request, v any, cc *cacheConfig) error {
	return p.sendRequestType(req, v, xmlHeader, cc)
}

func (p *Plex) sendRequestJSON(req *http.Request, v any, cc *cacheConfig) error {
	return p.sendRequestType(req, v, jsonHeader, cc)
}

func (p *Plex) doReq(req *http.Request) ([]byte, error) {
	res, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer dclose(res.Body)

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return content, err
	}
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return nil, errors.New(errRes.Message)
		}
		return nil, fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}
	return content, nil
}

func (p *Plex) sendRequestType(req *http.Request, v any, contentType string, cc *cacheConfig) error {
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", contentType)
	p.preprocessReq(req)

	var freshen bool
	content := []byte{}
	var key string
	if cc != nil {
		key = makeCacheKey(cc.prefix, *req)
	}
	if (cc == nil) || (cc.ttl == 0) {
		freshen = true
	} else {
		got, ok := p.cache.Get(key)
		if !ok {
			slog.Debug("cache not found, fetching fresh", "key", key)
			freshen = true
		} else {
			slog.Debug("using cache", "key", key)
			content = got.([]byte)
		}
	}

	if freshen {
		var err error
		content, err = p.doReq(req)
		if err != nil {
			return err
		}
		if (cc != nil) && (cc.ttl != 0) {
			p.cache.Set(key, content, cc.ttl)
		}
	}

	switch contentType {
	case "application/xml":
		if err := xml.NewDecoder(bytes.NewReader(content)).Decode(&v); err != nil && err != io.EOF {
			return err
		}
	case "application/json":
		if err := json.NewDecoder(bytes.NewReader(content)).Decode(&v); err != nil && err != io.EOF {
			return err
		}
	default:
		return errors.New("unknown content-type: " + contentType)
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

func (p *Plex) episodeID(show ShowTitle, season SeasonNumber, episode EpisodeNumber) (int, error) {
	shows, err := p.Shows.Match(show)
	if err != nil {
		return 0, err
	}
	episodes, err := p.Shows.Episodes(shows)
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
