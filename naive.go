package naive_cache

import (
	"fmt"
	"log"
	"sync"
)

// Getter load data for a key
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc implements Getter with a function.
type GetterFunc func(key string) ([]byte, error)

// Get impl the Getter interface
func (f GetterFunc) Get(key string) ([]byte, error)  {
	return f(key)
}

// Group is a cache namespace and associated data loaded spread over
type Group struct {
	name string
	getter Getter
	mainCache cache
}

var (
	mu sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup create a new instance of Group
func NewGroup(name string, getter Getter, cacheBytes int64) *Group {
	if getter == nil {
		 panic("getter func can not be ni")
	}

	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name:     name,
		getter:   getter,
		mainCache: cache{
			cacheBytes: cacheBytes,
		},
	}

	groups[name] = g
	return g
}

// GetGroup returns the named group previously created with NewGroup, or
// nil if there's no such group.
func  GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()

	return g
}

// Get value for a key from cache
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key can not be empty string")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Printf("[naive-cache] key:%s hint\n", key)
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	// user callback
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{b: cloneBytes(bytes)}
	g.mainCache.add(key, value)
	return value, nil
}