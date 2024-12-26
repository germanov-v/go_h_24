package hw04lrucache

func NewList() List {
	return new(list)
}

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	//List // Remove me after realization.
	// Place your code here.

	start     *ListItem
	end       *ListItem
	lenght    int
	items     []ListItem        // extra
	positions map[*ListItem]int // extra
}

func (l *list) Len() int {
	return l.lenght
}

func (l *list) Front() *ListItem {
	return l.start
}

func (l *list) Back() *ListItem {
	return l.end
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := ListItem{Value: v}
	item.Next = l.start
	l.start = &item
	l.lenght++
	return &item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := ListItem{Value: v}
	item.Prev = l.end
	l.end = &item
	l.lenght++
	return &item
}

func (l *list) Remove(i *ListItem) {

}

func (l *list) MoveToFront(i *ListItem) {

}
