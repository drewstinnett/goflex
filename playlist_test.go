package goflex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
