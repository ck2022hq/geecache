package geecache

import (
	"fmt"
	"geecache/singleflight"
	"log"
	"sync"
)

// A Getter loads data for a key.
type Getter interface {
	Get(key string) ([]byte, error)
}

// A GetterFunc implements Getter with a function.
type GetterFunc func(key string) ([]byte, error)

// Get implements Getter interface function
// 将函数类型转换成接口
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name     string
	getter   Getter
	peerPick PeerPicker
	cache    *Cache
	callOnce *singleflight.CallOnce
}

func NewGroup(name string, maxBytes int64, getter Getter, peerPick PeerPicker) *Group {
	return &Group{name, getter, peerPick, NewCache(maxBytes), singleflight.NewCallOnce()}
}

var (
	mtx    sync.RWMutex
	groups = make(map[string]*Group)
)

func AddGroup(name string, maxBytes int64, getter Getter, peerPick PeerPicker) *Group {
	mtx.Lock()
	defer mtx.Unlock()

	group := NewGroup(name, maxBytes, getter, peerPick)
	groups[name] = group
	return group
}

func GetGroup(name string) (*Group, bool) {
	mtx.RLock()
	defer mtx.RUnlock()

	g, ok := groups[name]
	return g, ok
}

func (g *Group) Get(key string) (ByteView, error) {
	if len(key) == 0 {
		log.Println("key cannot be empty")
		return ByteView{}, fmt.Errorf("key cannot be empty")
	}

	if g.peerPick == nil {
		return g.getLocal(key)
	}

	peer, ok := g.peerPick.PickPeer(key)
	if !ok {
		return g.getLocal(key)
	}

	val, err := g.callOnce.Call(key, func() (interface{}, error) {
		return g.loadRemote(peer, key)
	})
	return val.(ByteView), err
}

func (g *Group) getLocal(key string) (ByteView, error) {
	v, ok := g.cache.Get(key)
	if ok {
		log.Printf("[GeeCache] hit, group=%s key=%s", g.name, key)
		return v, nil
	}
	return g.loadLocal(key)
}

func (g *Group) loadRemote(peer PeerGetter, key string) (ByteView, error) {
	b, err := peer.Get(g.name, key)
	if err != nil {
		log.Printf("[GeeCache] load remote fail, group=%s key=%s err=%v", g.name, key, err)
		return ByteView{}, err
	}

	log.Printf("[GeeCache] load remote success, group=%s key=%s", g.name, key)
	return ByteView{cloneBytes(b)}, nil
}

func (g *Group) loadLocal(key string) (ByteView, error) {
	v, err := g.getter.Get(key)
	if err == nil {
		log.Printf("[GeeCache] load local success, group=%s key=%s", g.name, key)
		var b ByteView = ByteView{cloneBytes(v)}
		g.cache.Put(key, b)
		return b, nil
	}

	log.Printf("[GeeCache] load local fail, group=%s key=%s err=%v", g.name, key, err)
	return ByteView{}, err
}
