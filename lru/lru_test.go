package lru_test

import (
	"fmt"
	"geecache/lru"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func onEvicted(key string, value lru.Value) {
	fmt.Printf("[%s, %v] be evicted\n", key, value)
}

type mystring string

func (s mystring) Len() int {
	return len(s)
}

func TestLruCache_Get2(t *testing.T) {
	assert.Equal(t, 123, 123, "they should be equal")
	assert.NotEqual(t, 223, 123, "they should not be equal")

	cache := lru.NewLruCache(100, onEvicted)
	_, ok := cache.Get("haha")
	assert.False(t, ok, "no haha in cache")

	cache.Put("1", mystring(strings.Repeat("a", 20)))
	cache.Put("2", mystring(strings.Repeat("b", 20)))
	cache.Put("3", mystring(strings.Repeat("c", 20)))
	cache.Put("4", mystring(strings.Repeat("d", 20)))
	fmt.Println(cache)
	fmt.Printf("len:%d, size:%d\n", cache.Len(), cache.Bytes())
	cache.Put("5", mystring(strings.Repeat("e", 20)))
	fmt.Println(cache)

	val, ok := cache.Get("2")
	fmt.Println(cache)
	assert.True(t, ok)
	assert.Equal(t, "bbbbbbbbbbbbbbbbbbbb", string(val.(mystring)))
	cache.Put("6", mystring(strings.Repeat("f", 20)))
	fmt.Println(cache)
}

func TestLruCache_Get(t *testing.T) {
	type fields struct {
		maxBytes  int64
		OnEvicted func(key string, value lru.Value)
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   lru.Value
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := lru.NewLruCache(tt.fields.maxBytes, tt.fields.OnEvicted)
			got, got1 := cache.Get(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LruCache.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("LruCache.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
