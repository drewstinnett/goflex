package goflex

import (
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

func makeCacheKey(prefix string, req http.Request) string {
	var sb strings.Builder
	sb.WriteString(prefix)
	sb.WriteString(":")
	sb.WriteString(req.URL.String())
	return sb.String()
}

type cacheConfig struct {
	prefix string
	ttl    time.Duration
}

// cache stores data with TTL handling.
type cache struct {
	mutex      sync.RWMutex
	data       map[string]cacheItem
	expiries   map[string]time.Time
	gcRunning  bool
	gcInterval time.Duration
}

type cacheItem struct {
	value any
}

// newCache creates a new Cache.
func newCache() *cache {
	return &cache{
		data:       make(map[string]cacheItem),
		expiries:   make(map[string]time.Time),
		gcInterval: time.Minute * 1,
	}
}

// NewCacheWithGCI creates a new cache with a custom garbage collection interval.
func newCacheWithGC(i time.Duration) *cache {
	c := newCache()
	c.gcInterval = i
	return c
}

func (c *cache) exists(key string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	_, found := c.data[key]
	return found
}

func (c *cache) Get(key string) (any, bool) {
	c.mutex.RLock()
	item, found := c.data[key]
	expiry, expiryExists := c.expiries[key]
	isExpired := expiryExists && time.Now().After(expiry)
	c.mutex.RUnlock() // only unlock once

	if !found {
		return nil, false
	}
	if isExpired {
		c.Delete(key)
		return nil, false
	}
	return item.value, true
}

func (c *cache) Set(key string, value any, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = cacheItem{value: value}
	if ttl > 0 {
		c.expiries[key] = time.Now().Add(ttl)
	}

	if !c.gcRunning {
		c.gcRunning = true
		go c.startGC()
	}
}

func (c *cache) deleteWithKey(key string) {
	delete(c.data, key)
	delete(c.expiries, key)
}

func (c *cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.deleteWithKey(key)
}

func (c *cache) DeletePrefix(prefix string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for k := range c.data {
		if strings.HasPrefix(k, prefix) {
			c.deleteWithKey(k)
		}
	}
}

func (c *cache) startGC() {
	slog.Debug("starting garbage collection", "interval", c.gcInterval)
	ticker := time.NewTicker(c.gcInterval)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		for k, expiry := range c.expiries {
			if time.Now().After(expiry) {
				slog.Debug("garbage collecting cache key", "key", k, "expired", expiry)
				c.deleteWithKey(k)
			}
		}
		if len(c.expiries) == 0 {
			c.gcRunning = false
			c.mutex.Unlock()
			break
		}
		c.mutex.Unlock()
	}
}
