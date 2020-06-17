package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash function maps bytes to uint32
// User could provide their own hash function
type Hash func(data []byte) uint32

/**
hashring uses Map as an internal data structure to store keys.


*/
type Map struct {
	hash     Hash
	replicas int            // indicates the number of virtual nodes
	keys     []int          // keys sorted array
	hashMap  map[int]string // mapping between virtual nodes and actual keys
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))

	// find virtual node for this key
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	/**
	  idx%len(m.keys) == 0 , return m.keys[0] when idx == len(m.keys)
	  This makes they keys array a circular array

	  find the actual node for the given key
	*/
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
