package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	Cache // Remove me after realization.

	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lc *lruCache) Set(key Key, value interface{}) bool {

	return true
}

func (lc *lruCache) Clear() {
	lc = &lruCache{
		capacity: lc.capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, lc.capacity),
	}
}
