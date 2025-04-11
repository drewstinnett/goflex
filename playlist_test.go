package goflex

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlaylistList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch url := r.URL.Path; {
		case strings.HasSuffix(url, "/playlists"):
			expected, err := os.ReadFile("testdata/playlists.xml")
			require.NoError(t, err)
			fmt.Fprint(w, string(expected))
		case strings.HasSuffix(url, "/items"):
			expected, err := os.ReadFile("testdata/playlist-items.xml")
			require.NoError(t, err)
			fmt.Fprint(w, string(expected))
		default:
			panic("how did we get here? url: " + r.URL.Path)
		}
	}))
	// srv := srvFile(t, "testdata/playlists.xml")
	defer srv.Close()

	c, err := New(WithBaseURL(srv.URL), WithToken("test-token"))
	require.NoError(t, err)

	got, err := c.Playlists.List()
	require.NoError(t, err)

	require.Len(t, got, 7)

	ij, err := c.Playlists.GetWithName("Impractical Jokers (Randomized)")
	require.NoError(t, err)

	episodes, err := c.Playlists.Episodes(*ij)
	require.NoError(t, err)
	require.Len(t, episodes, 120)

	// fmt.Fprintf(os.Stderr, "ij: %v\n", episodes)

	eid, err := c.Playlists.EpisodeID(*ij, "Impractical Jokers", 8, 26)
	require.NoError(t, err)
	require.Equal(t, 32074, eid)
}

func TestNewRandomizeRequest(t *testing.T) {
	tests := []struct {
		name     string
		playlist PlaylistTitle
		series   []RandomizeSeries
		opts     []RandomizeRequestOpt
		wantErr  error
		wantReq  *RandomizeRequest
	}{
		{
			name:     "Valid request",
			playlist: "MyPlaylist",
			series: []RandomizeSeries{
				{
					Filter: EpisodeFilter{
						Show: "Show1",
					},
					LookbackDays: 7,
					// RefillAt: 10,
				},
			},
			opts:    nil,
			wantErr: nil,
			wantReq: &RandomizeRequest{
				Playlist: "MyPlaylist",
				Series: []RandomizeSeries{
					{
						Filter: EpisodeFilter{
							Show: "Show1",
						},
						LookbackDays: 7,
						// RefillAt: 10,
					},
				},
				RefillAt: 5,
			},
		},
		{
			name:     "Empty playlist",
			playlist: "",
			series: []RandomizeSeries{
				{
					Filter: EpisodeFilter{
						Show: "Show1",
					},
				},
			},
			opts:    nil,
			wantErr: errors.New("playlist must not be empty"),
			wantReq: nil,
		},
		{
			name:     "Empty series",
			playlist: "MyPlaylist",
			series:   []RandomizeSeries{},
			opts:     nil,
			wantErr:  errors.New("series muset not be empty"),
			wantReq:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := NewRandomizeRequest(tt.playlist, tt.series, tt.opts...)
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantReq, req)
			}
		})
	}
}
