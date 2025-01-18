package plexrando

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubtract(t *testing.T) {
	given := EpisodeList{
		{ID: 1},
		{ID: 2},
		{ID: 3},
	}
	remaining, removed := given.Subtract(EpisodeList{
		{ID: 2},
	})
	require.Equal(t, EpisodeList{
		{ID: 1},
		{ID: 3},
	}, remaining)
	require.Equal(t, EpisodeList{
		{ID: 2},
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
		got := everything.Seasons(tt.givenStart, tt.givenEnd)
		require.Equal(t, tt.expect, got, desc)
	}
}

func TestUpdateCache(t *testing.T) {
	pl := Playlist{}
	require.NotPanics(t, func() {
		pl.updateEpisodeCache(&EpisodeList{
			{Show: "foo", Season: 1, Episode: 1, PlaylistItemID: 5},
			{Show: "foo", Season: 1, Episode: 2, PlaylistItemID: 6},
		})
	})
	require.Equal(t, playlistEpisodeCache{
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
