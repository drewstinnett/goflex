package goflex

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalPlaylistItems(t *testing.T) {
	b, err := os.ReadFile("testdata/playlist-items.xml")
	require.NoError(t, err)
	v := PlaylistResponse{}
	require.NoError(t, xml.Unmarshal(b, &v))
	require.Equal(t, "Impractical Jokers (Randomized)", v.Title)
}
