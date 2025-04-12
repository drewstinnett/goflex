package goflex

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSortEpisodes(t *testing.T) {
	episodes := EpisodeList{
		Episode{Season: 1, Episode: 3},
		Episode{Season: 2, Episode: 2},
		Episode{Season: 1, Episode: 1},
	}
	sort.Sort(episodes)
	require.Equal(t, EpisodeList{
		Episode{Season: 1, Episode: 1},
		Episode{Season: 1, Episode: 3},
		Episode{Season: 2, Episode: 2},
	}, episodes)
}

func TestShowSeasons(t *testing.T) {
	/*
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {

			expected := SeasonMap{
				1: &Season{ID: 1, Index: 1},
				2: &Season{ID: 2, Index: 1},
			}
			expectedB, err := json.Marshal(expected)
			require.NoError(t, err)
			fmt.Fprint(w, expectedB)
		}))
	*/
	svr := srvFile(t, "testdata/seasons.json")
	defer svr.Close()

	p, err := New(
		WithBaseURL(svr.URL),
		WithHTTPClient(http.DefaultClient),
		WithToken("test-token"),
	)
	require.NoError(t, err)
	seasons, err := p.Shows.Seasons(Show{ID: 2, Title: "Fake Show"})
	require.NoError(t, err)
	require.NotNil(t, seasons)
	require.Len(t, *seasons, 21)
	s := *seasons
	assert.EqualValues(t, 1, s[1].Index)
	assert.EqualValues(t, 20, s[20].Index)
}

func TestEpisodeMapList(t *testing.T) {
	require.Equal(t,
		EpisodeList{
			Episode{Season: 1, Episode: 1},
			Episode{Season: 1, Episode: 2},
			Episode{Season: 1, Episode: 3},
		},
		EpisodeMap{
			2: &Episode{Season: 1, Episode: 2},
			3: &Episode{Season: 1, Episode: 3},
			1: &Episode{Season: 1, Episode: 1},
		}.List(),
	)
}

func TestMatchShows(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(os.Stderr, "REQ: %+v\n", req)
		if strings.HasSuffix(req.URL.String(), "/library/sections/") {
			expected, err := os.ReadFile("./testdata/libraries.xml")
			assert.NoError(t, err)
			fmt.Fprint(w, string(expected))
			return
		}
		expected, err := os.ReadFile("./testdata/shows.xml")
		assert.NoError(t, err)
		fmt.Fprint(w, string(expected))
	}))

	p, err := New(WithBaseURL(svr.URL), WithToken("test-token"))
	require.NoError(t, err)

	got, err := p.Shows.Match("American Dad!")
	require.NoError(t, err)
	require.Len(t, got, 2)
	// 2 versions of the show because it's in multiple libraries
	require.EqualValues(t, "American Dad!", got[0].Title)
	require.EqualValues(t, "American Dad!", got[1].Title)
}

func TestShowExists(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(os.Stderr, "REQ: %+v\n", req)
		if strings.HasSuffix(req.URL.String(), "/library/sections/") {
			expected, err := os.ReadFile("./testdata/libraries.xml")
			assert.NoError(t, err)
			fmt.Fprint(w, string(expected))
			return
		}
		expected, err := os.ReadFile("./testdata/shows.xml")
		assert.NoError(t, err)
		fmt.Fprint(w, string(expected))
	}))

	p, err := New(WithBaseURL(svr.URL), WithToken("test-token"))
	require.NoError(t, err)
	got, err := p.Shows.Exists("never-exists")
	require.NoError(t, err)
	require.False(t, got)

	got, err = p.Shows.Exists("American Dad!")
	require.NoError(t, err)
	require.True(t, got)
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
		givenStart SeasonNumber
		givenEnd   SeasonNumber
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

func TestSeasonMap(t *testing.T) {
	require.Equal(t, SeasonList{
		{Index: 1},
		{Index: 2},
		{Index: 3},
	}, SeasonMap{
		2: &Season{Index: 2},
		1: &Season{Index: 1},
		3: &Season{Index: 3},
	}.sorted())
}
