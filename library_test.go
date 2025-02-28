package goflex

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLibraries(t *testing.T) {
	expected, err := os.ReadFile("./testdata/libraries.xml")
	require.NoError(t, err)
	hits := 0
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hits++
		fmt.Fprint(w, string(expected))
	}))
	defer svr.Close()

	p, err := New(
		WithBaseURL(svr.URL),
		WithHTTPClient(http.DefaultClient),
		WithToken("test-token"),
	)
	require.NoError(t, err)
	got, err := p.Library.List()
	require.NoError(t, err)
	require.Equal(t, 6, len(got))
	require.Equal(t, 1, hits)

	// Call again, but this time we should use the cache, so no additional hits
	_, err = p.Library.List()
	require.NoError(t, err)
	require.Equal(t, 1, hits)
}
