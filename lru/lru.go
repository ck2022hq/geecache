package lru

import (
	"bytes"
	"container/list"
	"fmt"
)

type LruCache struct {
	values map[string]*list.Element
	// list.Element.value type: *entry
	recentList *list.List

	maxBytes int64
	nbytes   int64

	onEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

func (e *entry) Len() int {
	return len(e.key) + e.value.Len()
}

func NewLruCache(maxBytes int64, onEvicted func(key string, value Value)) *LruCache {
	return &LruCache{
		values:     make(map[string]*list.Element),
		recentList: list.New(),
		nbytes:     0,
		maxBytes:   maxBytes,
		onEvicted:  onEvicted,
	}
}

func (cache *LruCache) get(key string) (*entry, bool) {
	elem, ok := cache.values[key]
	if !ok {
		return nil, false
	}

	cache.recentList.MoveToFront(elem)
	kv := elem.Value.(*entry)
	return kv, true
}

func (cache *LruCache) Get(key string) (Value, bool) {
	kv, ok := cache.get(key)
	if !ok {
		return nil, ok
	}

	return kv.value, ok
}

func (cache *LruCache) removeOldest() {
	if len(cache.values) == 0 {
		return
	}

	elem := cache.recentList.Back()
	kv := elem.Value.(*entry)
	cache.nbytes -= int64(kv.Len())

	delete(cache.values, kv.key)
	cache.recentList.Remove(elem)

	if cache.onEvicted != nil {
		cache.onEvicted(kv.key, kv.value)
	}
}

func (cache *LruCache) Put(key string, value Value) {
	kv, ok := cache.get(key)
	if ok {
		cache.nbytes += int64(value.Len() - kv.value.Len())
		kv.value = value
	} else {
		e := &entry{key, value}
		cache.recentList.PushFront(e)
		cache.values[key] = cache.recentList.Front()
		cache.nbytes += int64(e.Len())
	}

	for cache.nbytes > cache.maxBytes {
		cache.removeOldest()
	}
}

func (cache *LruCache) Len() int {
	return cache.recentList.Len()
}

func (cache *LruCache) Bytes() int64 {
	return cache.nbytes
}

func (cache *LruCache) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("{")
	for elem := cache.recentList.Front(); elem != nil; elem = elem.Next() {
		kv := elem.Value.(*entry)
		buffer.WriteString(fmt.Sprintf("[%s, %v],", kv.key, kv.value))
	}
	buffer.WriteString("}")
	return buffer.String()
}
