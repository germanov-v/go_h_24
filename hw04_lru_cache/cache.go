package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type itemCache struct {
	value interface{}
	key   Key
}

type lruCache struct {
	//Cache    // Remove me after realization.
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
		item.Value = itemCache{value: value, key: key}
		lc.queue.MoveToFront(item)
		return true
	}
	if lc.capacity == lc.queue.Len() {
		removeItem := lc.queue.Back()
		removeItemValue := removeItem.Value.(itemCache)
		lc.queue.Remove(removeItem)
		delete(lc.items, removeItemValue.key)
	}
	newItem := lc.queue.PushFront(itemCache{value: value, key: key})
	lc.items[key] = newItem

	return false
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	item, ok := lc.items[key]
	if !ok {
		return nil, ok
	}
	lc.queue.MoveToFront(item)
	needValue := item.Value.(itemCache)
	return needValue.value, ok
}

func (lc *lruCache) Clear() {
	lc.queue = NewList()
	lc.items = make(map[Key]*ListItem, lc.capacity)
}
