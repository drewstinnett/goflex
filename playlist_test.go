package goflex

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

/*
func TestUpdateCache(t *testing.T) {
	pl := Playlist{}
	require.NotPanics(t, func() {
		pl.updateEpisodeCache(&EpisodeList{
			{Show: "foo", Season: 1, Episode: 1, PlaylistItemID: 5},
			{Show: "foo", Season: 1, Episode: 2, PlaylistItemID: 6},
		})
	})
	require.Equal(t, PlaylistEpisodeCache{
		"foo": {
			1: {
				1: 5,
				2: 6,
			},
		},
	}, pl.episodeCache)

	k, err := pl.episodeKey("foo", 1, 2)
	require.NoError(t, err)
	require.Equal(t, 6, k)
}
*/

func TestNewRandomizeRequest(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		t.Errorf("got request I shouldn't have: %v", req.URL)
	}))
	defer svr.Close()

	p, err := New(
		WithBaseURL(svr.URL),
		WithToken("test-token"),
	)
	require.NoError(t, err)
	require.NotNil(t, p)
	// p.Playlists.episodeCache = nil

	r, err := NewRandomizeRequest("test-playlist", []RandomizeSeries{
		{
			Filter:   EpisodeFilter{Show: "show1"},
			Lookback: toPTR(time.Hour * 24 * 30),
			RefillAt: 1,
		},
	})
	require.NoError(t, err)
	require.NotNil(t, r)

	got, err := p.Playlists.Randomize(*r)
	require.NoError(t, err)
	require.NotNil(t, got)
}
