package cacher

import "sync"

type Key string

type Value interface{}

type Cache interface {
	Get(key Key) (Value, bool)
	Set(key Key, value Value) bool
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	mu       sync.RWMutex
	items    map[Key]*ListItem
}

func New(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Get(key Key) (Value, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(item)
		return item.Value, true
	}
	return nil, false
}

func (c *lruCache) Set(key Key, value Value) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, wasInCache := c.items[key]

	if wasInCache {
		item.Value = value
		c.queue.MoveToFront(item)
	} else {
		item = &ListItem{Value: value}
		c.items[key] = c.queue.PushFront(key, item.Value)

		if c.queue.Len() > c.capacity {
			back := c.queue.Back()
			c.queue.Remove(back)
			delete(c.items, back.Key)
		}
	}

	return wasInCache
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[Key]*ListItem)
}
