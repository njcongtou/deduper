package lru

import (
	"container/list"
)

/*
   in memory LRU cache.
*/
type Cache struct {
	// max memory usage in bytes.
	maxBytes int64

	// current memory usage in bytes.
	currentBytes int64

	// a doubly Linked List.
	// fast element removal, moving to front or deletion.
	dList *list.List

	// map : key is string, and value is Linked list element.
	cache map[string]*list.Element

	// a function executed on an entry is purged.
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value // actual value of the key
}

type Value interface {
	Len() int // size of memory usage of a key in bytes
}

// instantiate a cache object
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:     maxBytes,
		currentBytes: 0,
		dList:        list.New(),
		cache:        make(map[string]*list.Element),
		OnEvicted:    onEvicted, // a function is executed when this cache is deleted.
	}
}

/**
  Move recent accessed element to the front of the List
*/
func (c *Cache) Get(key string) (value Value, ok bool) {
	if element, ok := c.cache[key]; ok {
		c.dList.MoveToFront(element)
		kv := element.Value.(*entry)
		return kv.value, true
	}
	return
}

/**
  LRU policy: remove the last element in the list.
*/
func (c *Cache) DeleteOldest() {
	element := c.dList.Back() // find the last element in the list.
	if element != nil {
		c.dList.Remove(element)
		kv := element.Value.(*entry)
		delete(c.cache, kv.key)
		c.currentBytes -= int64(len(kv.key) + kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if element, ok := c.cache[key]; ok {
		// this is an existing key,
		// so add it into the doubly list.
		// update the size, delta is: new value size - old value size
		c.dList.MoveToFront(element)
		kv := element.Value.(*entry)
		c.currentBytes += int64(value.Len() - kv.value.Len())
		kv.value = value
	} else {
		// this is a new key
		element = c.dList.PushFront(&entry{key, value})
		c.cache[key] = element
		c.currentBytes += int64(value.Len() + len(key))
	}

	// remove the last element if current memory usage exceeds the limit
	// for loop to keep deleting last elements until limit is not exceeded.
	for c.maxBytes != 0 && c.maxBytes < c.currentBytes {
		c.DeleteOldest()
	}
}

// the number of cache entries
func (c *Cache) Len() int {
	return c.dList.Len()
}
