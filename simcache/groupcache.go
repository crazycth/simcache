package simcache

import (
	"fmt"
	"log"
	"sync"

	"github.com/crazycth/simcache/singleflight"
	"github.com/crazycth/simcache/tools/async"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	getter    Getter //e.g. seach database
	mainCache cache
	peers     PeerPicker
	loader    *singleflight.Group //avoid Cache breakdown
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

//NewGroup create a new instance of Group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter is invalid!")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:   name,
		getter: getter,
		mainCache: cache{
			cacheBytes: cacheBytes,
		},
		loader: &singleflight.Group{
			FutureM: make(map[string]*async.Future),
			CallM:   make(map[string]*singleflight.Call),
		},
	}
	groups[name] = g
	return g
}

//GetGroup returns the named group previously created with NewGroup
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	g := groups[name]
	return g
}

func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("empty key")
	}
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[simcache][Get] hit cache , key: ", key)
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
			log.Println("[simcache][load] Failed to get from peer", err.Error())
		}
	}
	return g.getLocally(key)
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytesI, err := g.loader.Query(key, func() (interface{}, error) {
		bytes, err := g.getter.Get(key)
		if err != nil {
			log.Println("[simcache][getLocally] getter error with ", err.Error())
			return ByteView{}, err
		}
		return bytes, err
	})
	if err != nil {
		log.Println("[simcache][getLocally] getter error with ", err.Error())
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytesI.([]byte))}
	g.mainCache.add(key, value)
	return value, nil
}

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

func (g *Group) QueryPeer(name string, key string) (ByteView, error) {
	if g.peers == nil {
		return ByteView{}, fmt.Errorf("server peers nil error")
	}

	if peer, ok := g.peers.PickPeer(key); ok {
		return g.getFromPeer(peer, key)
	}

	return ByteView{}, fmt.Errorf("server pick peer error")
}
