package geecache

import (
	"geecache/lru"
	"sync"
)

type Cache struct {
	sync.Mutex

	lru *lru.LruCache
}

func NewCache(maxBytes int64) *Cache {
	return &Cache{
		lru: lru.NewLruCache(maxBytes, nil),
	}
}

func (cache *Cache) Put(key string, value ByteView) {
	cache.Lock()
	defer cache.Unlock()

	cache.lru.Put(key, value)
}

func (cache *Cache) Get(key string) (ByteView, bool) {
	cache.Lock()
	defer cache.Unlock()

	if cache.lru == nil {
		return ByteView{}, false
	}

	v, ok := cache.lru.Get(key)
	if !ok {
		return ByteView{}, ok
	}
	return v.(ByteView), ok
}
