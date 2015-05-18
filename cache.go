package composer

import (
	"time"
)

type Cache interface {
	Get(key string) (string, bool)
	Set(key string, value string, expires time.Time)
}

type cacheItem struct {
	value   string
	expires time.Time
}

type MemoryCache struct {
	items map[string]cacheItem
}

func NewMemoryCache() Cache {
	return &MemoryCache{items: make(map[string]cacheItem)}
}

func (c *MemoryCache) Get(key string) (string, bool) {
	item, ok := c.items[key]

	if ok {
		now := time.Now()

		if now.Before(item.expires) {
			return item.value, true
		}

		delete(c.items, key)
	}

	return "", false
}

func (c *MemoryCache) Set(key string, value string, expires time.Time) {
	item := cacheItem{value: value, expires: expires}
	c.items[key] = item
}
