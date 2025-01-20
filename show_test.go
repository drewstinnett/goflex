package goflex

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatchShows(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(os.Stderr, "REQ: %+v\n", req)
		if strings.HasSuffix(req.URL.String(), "/library/sections/") {
			expected, err := os.ReadFile("./testdata/libraries.xml")
			require.NoError(t, err)
			fmt.Fprint(w, string(expected))
			return
		}
		expected, err := os.ReadFile("./testdata/shows.xml")
		require.NoError(t, err)
		fmt.Fprint(w, string(expected))
	}))

	p, err := New(WithBaseURL(svr.URL), WithToken("test-token"))
	require.NoError(t, err)

	got, err := p.Shows.Match("American Dad!")
	require.NoError(t, err)
	require.Equal(t, 2, len(got))
	// 2 versions of the show because it's in multiple libraries
	require.Equal(t, "American Dad!", got[0].Title)
	require.Equal(t, "American Dad!", got[1].Title)
}

func TestShowExists(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(os.Stderr, "REQ: %+v\n", req)
		if strings.HasSuffix(req.URL.String(), "/library/sections/") {
			expected, err := os.ReadFile("./testdata/libraries.xml")
			require.NoError(t, err)
			fmt.Fprint(w, string(expected))
			return
		}
		expected, err := os.ReadFile("./testdata/shows.xml")
		require.NoError(t, err)
		fmt.Fprint(w, string(expected))
	}))

	p, err := New(WithBaseURL(svr.URL), WithToken("test-token"))
	require.NoError(t, err)
	got, err := p.Shows.Exists("never-exists")
	require.NoError(t, err)
	require.Equal(t, false, got)

	got, err = p.Shows.Exists("American Dad!")
	require.NoError(t, err)
	require.Equal(t, true, got)
}

func TestSubtract(t *testing.T) {
	given := EpisodeList{
		{Show: "foo", Season: 1, Episode: 1},
		{Show: "foo", Season: 1, Episode: 2},
		{Show: "foo", Season: 1, Episode: 3},
		{Show: "foo", Season: 1, Episode: 3},
	}
	remaining, removed := given.Subtract(EpisodeList{
		{Show: "foo", Season: 1, Episode: 2},
	})
	require.Equal(t, EpisodeList{
		{Show: "foo", Season: 1, Episode: 1},
		{Show: "foo", Season: 1, Episode: 3},
	}, remaining)
	require.Equal(t, EpisodeList{
		{Show: "foo", Season: 1, Episode: 2},
	}, removed)
}

func TestEpisodeSeasons(t *testing.T) {
	everything := EpisodeList{
		{Title: "s01e01", Season: 1},
		{Title: "s02e01", Season: 2},
		{Title: "s03e01", Season: 3},
	}
	tests := map[string]struct {
		givenStart int
		givenEnd   int
		expect     EpisodeList
	}{
		"start-and-end": {
			givenStart: 1,
			givenEnd:   2,
			expect: EpisodeList{
				{Title: "s01e01", Season: 1},
				{Title: "s02e01", Season: 2},
			},
		},
		"start-only": {
			givenStart: 2,
			expect: EpisodeList{
				{Title: "s02e01", Season: 2},
				{Title: "s03e01", Season: 3},
			},
		},
		"end-only": {
			givenEnd: 2,
			expect: EpisodeList{
				{Title: "s01e01", Season: 1},
				{Title: "s02e01", Season: 2},
			},
		},
	}

	for desc, tt := range tests {
		require.Equal(t, tt.expect, everything.Seasons(tt.givenStart, tt.givenEnd), desc)
	}
}
