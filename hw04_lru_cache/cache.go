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

	if item, ok := lc.items[key]; ok {
		lc.queue.MoveToFront(item)
	} else {
		if lc.capacity < lc.queue.Len() {
			lc.queue.Remove(lc.queue.Back())
		}
		newItem := lc.queue.PushFront(item)
		lc.items[key] = newItem
	}

	return true
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	item, ok := lc.items[key]
	return item.Value, ok
}

func (lc *lruCache) Clear() {
	lc = &lruCache{
		capacity: lc.capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, lc.capacity),
	}
}
