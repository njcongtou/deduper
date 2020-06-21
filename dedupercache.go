package deduper

import (
	"fmt"
	"log"
	"sync"
)

// loads data for a key.
type Getter interface {
	Get(key string) ([]byte, error)
}

// A GetterFunc implements Getter with a function.
// If key is not present in the cache, callback will try to get it from local db.
type GetterFunc func(key string) ([]byte, error)

// implements Get function defined in Getter interface
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

/** Group definition as follows **/

// Group is a cache namespace and associated loaded data
type Group struct {
	name      string
	getter    Getter
	mainCache cache
	peers     PeerPicker
}

var (
	mu sync.RWMutex
	/**
	Groups is major data structure that holdes a collected of cached values. such as
	Students' score, students' info, students' courses.
	*/
	groups = make(map[string]*Group)
)

// NewGroup creates a new instance of Group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}

	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

// GetGroup returns the named group previously created with NewGroup, or
// nil if there's no such group
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	g := groups[name]
	return g
}

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	fmt.Println("RegisterPeers in progress, old peers:")
	fmt.Println(g.peers)
	g.peers = peers
	fmt.Println("RegisterPeers completed, new peers:")
	fmt.Println(g.peers)
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[dedupercache] hit")
		return v, nil
	}
	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	if g.peers != nil {
		if peer, ok := g.peers.PickPeer(key); ok {
			if value, err = g.getFromPeer(peer, key); err == nil {
				return value, nil
			}
			log.Println("[dedupercache] Failed to get from peer", err)
		}
	}
	return g.getLocally(key)
}

// Get key from a remote peer
// TODO modify to get keys from 3 peers
func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	log.Println("[dedupercache] populated key from local db: ", key)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
