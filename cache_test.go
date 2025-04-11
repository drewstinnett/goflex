package goflex

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	c := NewCacheWithGC(time.Millisecond * 10)

	// Set a new key with a short cache
	c.Set("foo", "foo", time.Second*2)
	got, found := c.Get("foo")
	require.Equal(t, "foo", got)
	require.True(t, found)

	// Wait for cache to expire and make sure the value is gone
	time.Sleep(time.Second * 2)
	_, found = c.Get("foo")
	require.False(t, found)

	// Set and delete an item
	c.Set("bar", "bar", time.Hour*12)
	c.Delete("bar")
	_, found = c.Get("bar")
	require.False(t, found)

	// Set some prefix keys and then delete them
	c.Set("prefix:foo", "foo", time.Hour*12)
	c.Set("prefix:bar", "bar", time.Hour*12)
	c.Set("not-a-prefix:baz", "baz", time.Hour*12)
	c.DeletePrefix("prefix:")
	require.False(t, c.exists("prefix:foo"))
	require.False(t, c.exists("prefix:bar"))
	require.True(t, c.exists("not-a-prefix:baz"))
}
