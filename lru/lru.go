package lru

import "container/list"

// Cache is a LRU strategy cache
// not thread safe now
type Cache struct {
	// 0 represent no limit
	maxBytes int64
	// how many bytes we taked now
	nBytes int64
	// double-list internal in go
	ll    *list.List
	cache map[string]*list.Element
	// optional
	// when entry purged it will execute
	// like callback fnuc in `JavaScript`
	OnEvicted func(key string, value Value)
}

// entry need to keep key
// it is convient to delete in map when delete entry
// see leetcode :D
type entry struct {
	key   string
	value Value
}

// Value use len func to caclue hot many bytes it takes in mem
type Value interface {
	Len() int
}

// New is a constructor of `Cache`
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		nBytes:    0,
		ll:        list.New(),
		cache:     map[string]*list.Element{},
		OnEvicted: onEvicted,
	}
}

// Get will look element and move this node to list end
func (c *Cache) Get(key string) (value Value, ok bool) {
	if element, ok := c.cache[key]; ok {
		// lru feature
		c.ll.MoveToBack(element)
		kv := element.Value.(*entry)
		return kv.value, true
	}
	return
}

// Delete will delete the entry by lru
func (c *Cache) RemoveOldest() {
	element := c.ll.Front()
	if element != nil {
		c.ll.Remove(element)
		kv := element.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key) + kv.value.Len())

		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add will add an new entry in cache
func (c *Cache) Add(key string, value Value) {
	if element, ok := c.cache[key]; ok {
		// already exist
		c.ll.MoveToBack(element)
		kv := element.Value.(*entry)
		kv.value = value
		c.nBytes += int64(value.Len() - kv.value.Len())
	} else {
		element := c.ll.PushBack(&entry{key: key, value: value})
		c.cache[key] = element
		c.nBytes += int64(len(key) + value.Len())
	}

	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

// Len reutrn the number of elements
func (c *Cache) Len() int {
	return c.ll.Len()
}
