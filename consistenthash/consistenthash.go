package consistenthash

import (
	"fmt"
	"hash/crc32"
	"sort"
)

type Hash func(data []byte) uint32

type ConsistentHash struct {
	hash     Hash
	replicas int
	keys     []int // Sorted
	hashMap  map[int]string
}

func NewConsistentHash(replicas int, hash Hash) *ConsistentHash {
	if hash == nil {
		hash = crc32.ChecksumIEEE
	}
	return &ConsistentHash{replicas: replicas, hash: hash, hashMap: make(map[int]string)}
}

func (h *ConsistentHash) Add(hosts ...string) {
	for _, host := range hosts {
		for i := 0; i < h.replicas; i++ {
			key := fmt.Sprintf("%.4d%s", i, host)
			r := int(h.hash([]byte(key)))
			h.keys = append(h.keys, r)
			h.hashMap[r] = host
		}
	}
	sort.Ints(h.keys)
}

func (h *ConsistentHash) SelectHost(key string) string {
	r := int(h.hash([]byte(key)))

	idx := sort.Search(len(h.keys), func(i int) bool {
		return h.keys[i] >= r
	})
	return h.hashMap[h.keys[idx%len(h.keys)]]
}
