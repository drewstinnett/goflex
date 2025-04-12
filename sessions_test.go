package goflex

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestHistorySession(t *testing.T) {
	expected, err := os.ReadFile("./testdata/history-sessions.xml")
	require.NoError(t, err)
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, string(expected))
	}))
	defer svr.Close()

	p, err := New(
		WithBaseURL(svr.URL),
		WithHTTPClient(http.DefaultClient),
		WithToken("test-token"),
	)
	require.NoError(t, err)

	fakeNow := time.Date(2025, time.January, 10, 12, 0, 0, 0, time.Local)

	_, err = p.Sessions.HistoryEpisodes(fakeNow)
	require.NoError(t, err)
	// assert.Equal(t, 92, len(got))

	_, err = p.Sessions.HistoryEpisodes(fakeNow, "American Dad!", "Impractical Jokers")
	require.NoError(t, err)
	// assert.Equal(t, 65, len(got))
}
