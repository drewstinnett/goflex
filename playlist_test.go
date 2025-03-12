package goflex

import (
	"testing"
)

func TestUpdateCache(t *testing.T) {
	/*
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
	*/
}

/*
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
	p.Playlists.InitCache(map[string]*Playlist{
		"test-playlist": {
			Title: "test-playlist",
			ID:    1,
		},
	})
	p.Playlists.InitEpisodeCache(map[string]EpisodeList{
		"show1": {
			{Show: "show1", Season: 1, Episode: 1},
			{Show: "show1", Season: 1, Episode: 2},
			{Show: "show1", Season: 1, Episode: 3},
			{Show: "show1", Season: 1, Episode: 4},
			{Show: "show1", Season: 1, Episode: 5},
			{Show: "show1", Season: 1, Episode: 6},
			{Show: "show1", Season: 1, Episode: 7},
			{Show: "show1", Season: 1, Episode: 8},
		},
	})
	p.Sessions.InitHistory(EpisodeList{
		Episode{
			Show: "show1", Season: 1, Episode: 1, Watched: daysAgo(3),
		},
	})

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

func daysAgo(d int) *time.Time {
	return toPTR(time.Now().Add(-time.Hour * 24 * time.Duration(d)))
}
*/
