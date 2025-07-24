package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

var mutex sync.Mutex

func (c *lruCache) Set(key Key, value interface{}) bool {
	mutex.Lock()
	defer mutex.Unlock()

	item, isExist := c.items[key]

	if isExist {
		item.Value = cacheItem{key, value}
		c.queue.MoveToFront(item)
	} else {
		c.queue.PushFront(cacheItem{key, value})
		c.items[key] = c.queue.Front()

		if c.queue.Len() > c.capacity {
			delItem := c.queue.Back()
			c.queue.Remove(delItem)
			delete(c.items, delItem.Value.(cacheItem).key)
		}
	}

	return isExist
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	item, isExist := c.items[key]
	if isExist {
		c.queue.MoveToFront(item)

		return item.Value.(cacheItem).value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	mutex.Lock()
	defer mutex.Unlock()

	c.items = make(map[Key]*ListItem, c.capacity)
	c.queue = NewList()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
