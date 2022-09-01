package geecache

import (
	"fmt"
	"geecache/consistenthash"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// /cache?group=xxx&name=yyy
const maxBytes = 2 << 10
const defaultReplicas = 50

var DefaultGetter Getter = nil

const defaultBasePath = "/geecache/"

// HTTPPool implements PeerPicker for a pool of HTTP peers.
type HTTPPool struct {
	// this peer's base URL, e.g. "https://example.net:8000"
	self     string
	basePath string

	mu          sync.Mutex // guards peers and httpGetters
	peers       *consistenthash.ConsistentHash
	httpGetters map[string]*httpGetter // keyed by e.g. "http://10.0.0.2:8008"
}

var Pool *HTTPPool

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (pool *HTTPPool) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	pool.Log("%s %s", req.Method, req.URL.Path)

	groupname := req.URL.Query().Get("group")
	key := req.URL.Query().Get("key")
	if len(groupname) == 0 || len(key) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "groupname or key cannot be empty")
		return
	}

	val, err := pool.Process(groupname, key)

	if err != nil {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q in group %q\n", key, groupname)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(val.ByteSlice())
		// val.WriteTo(w)
	}
}

func (pool *HTTPPool) Process(groupname, key string) (ByteView, error) {
	group, ok := GetGroup(groupname)
	if !ok {
		group = AddGroup(groupname, maxBytes, DefaultGetter, pool)
	}

	return group.Get(key)
}

type httpGetter struct {
	baseURL string
}

func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf(
		"%v?group=%v&key=%v",
		strings.TrimRight(h.baseURL, "/"),
		url.QueryEscape(group),
		url.QueryEscape(key),
	)
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Printf("http get %s fail, error=%v", u, res.Status)
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	return bytes, nil
}

// Set updates the pool's list of peers.
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.peers = consistenthash.NewConsistentHash(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}

func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if peer := p.peers.SelectHost(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}
