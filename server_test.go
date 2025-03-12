package goflex

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func srvFile(t *testing.T, f string) *httptest.Server {
	expected, err := os.ReadFile(f)
	require.NoError(t, err)
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, string(expected))
	}))
}

func TestSearch(t *testing.T) {
	svr := srvFile(t, "./testdata/search-results.json")
	defer svr.Close()

	p, err := New(
		WithBaseURL(svr.URL),
		WithHTTPClient(http.DefaultClient),
		WithToken("test-token"),
	)
	require.NoError(t, err)
	got, err := p.Server.Search("family")
	require.NoError(t, err)

	shows, err := got.Shows()
	require.NoError(t, err)
	require.Equal(t, 1, len(shows))
	require.EqualValues(t, "Family Guy", shows[0].Title)

	episodes, err := got.Episodes()
	require.NoError(t, err)
	require.EqualValues(t, 15, len(episodes))
	require.EqualValues(t, "'Family Guy' Through The Years", episodes[0].Title)
}
