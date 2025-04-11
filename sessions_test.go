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

	got, err := p.Sessions.HistoryEpisodes(toPTR(time.Now()))
	require.NoError(t, err)
	require.Equal(t, 92, len(got))

	got, err = p.Sessions.HistoryEpisodes(toPTR(time.Now()), "American Dad!", "Impractical Jokers")
	require.NoError(t, err)
	require.Equal(t, 65, len(got))
}
