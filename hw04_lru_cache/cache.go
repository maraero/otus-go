package hw04lrucache

import "fmt"

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

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) addItem(key Key, value interface{}) {
	item := cacheItem{key: key, value: value}
	listItemPtr := c.queue.PushFront(item)
	c.items[key] = listItemPtr
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if item, ok := c.items[key]; ok {
		c.queue.Remove(item)
		c.addItem(key, value)
		return true
	}

	c.addItem(key, value)

	if c.queue.Len() > c.capacity {
		lastItem := c.queue.Back()
		c.queue.Remove(lastItem)
		
		if valWithKey, ok := lastItem.Value.(cacheItem); ok {
			delete(c.items, valWithKey.key)
		}
	}

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	if item, ok := c.items[key]; ok {
		fmt.Println(item, c.queue)
		c.queue.MoveToFront(item)

		if valWithKey, ok := item.Value.(cacheItem); ok {
			return valWithKey.value, true
		}
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.items = make(map[Key]*ListItem, c.capacity)
	c.queue = NewList()
}