package plexrando

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubtract(t *testing.T) {
	given := EpisodeList{
		{ID: "foo"},
		{ID: "bar"},
		{ID: "baz"},
	}
	remaining, removed := given.Subtract(EpisodeList{
		{ID: "bar"},
	})
	require.Equal(t, EpisodeList{
		{ID: "foo"},
		{ID: "baz"},
	}, remaining)
	require.Equal(t, EpisodeList{
		{ID: "bar"},
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
