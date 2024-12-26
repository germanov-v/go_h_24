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
}

func (l *list) Len() int {
	return 1
}

func (l *list) Front() *ListItem {
	return new(ListItem)
}

func (l *list) Back() *ListItem {
	return new(ListItem)
}

func (l *list) PushFront(v interface{}) *ListItem {
	return new(ListItem)
}

func (l *list) PushBack(v interface{}) *ListItem {
	return new(ListItem)
}

func (l *list) Remove(i *ListItem) {

}

func (l *list) MoveToFront(i *ListItem) {

}
