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

	start  *ListItem
	end    *ListItem
	lenght int
	//items     []ListItem        // extra
	//positions map[*ListItem]int // extra
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
	if l.start != nil {
		item.Next = l.start
		l.start.Prev = &item
	} else {
		l.end = &item
	}
	l.start = &item
	l.lenght++
	return &item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := ListItem{Value: v}
	if l.end != nil {
		item.Prev = l.end
		l.end.Next = &item
	} else {
		l.start = &item
	}

	l.end = &item
	l.lenght++
	return &item
}

// TODO: has param 'i' actual pointers ? => O(1) without loop
func (l *list) Remove(i *ListItem) {

	if i == nil {
		return
	}

	if i.Prev != nil && i.Next != nil {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
		i.Prev = nil
	} else if i.Prev != nil {
		i.Prev.Next = nil
		l.end = i.Prev
	} else { // i.Next != nil
		i.Next.Prev = nil
		l.start = i.Next

	}
	i.Next = nil
	i.Prev = nil
	l.lenght--
}

func (l *list) MoveToFront(i *ListItem) {

	if l.start == i && i.Prev != nil {
		panic("plan went wrong: l.start == i && i.Prev != nil")
	} else if l.start == i {
		return
	}

	if i.Prev != nil && i.Next != nil {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	} else { //  i.Prev != nil
		i.Prev.Next = nil
		l.end = i.Prev
	}
	i.Next = l.start
	l.start.Prev = i
	i.Prev = nil
	l.start = i
}
