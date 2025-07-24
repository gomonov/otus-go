package hw04lrucache

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
	length    int
	firstItem *ListItem
	lastItem  *ListItem
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.firstItem
}

func (l *list) Back() *ListItem {
	return l.lastItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v, Next: l.firstItem, Prev: nil}
	if l.length == 0 {
		l.lastItem = item
	} else {
		l.firstItem.Prev = item
	}

	l.firstItem = item
	l.length++

	return l.firstItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v, Next: nil, Prev: l.lastItem}
	if l.length == 0 {
		l.firstItem = item
	} else {
		l.lastItem.Next = item
	}
	l.lastItem = item
	l.length++

	return l.lastItem
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.firstItem = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.lastItem = i.Prev
	}

	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}
