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
	front *ListItem
	back *ListItem
	len int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	elem := &ListItem{Value: v}

	switch l.Len() {
	case 0:
		l.front = elem
		l.back = elem
	case 1:
		l.front = elem
		l.front.Next = l.back
		l.back.Prev = l.front
	default:
		elem.Next = l.front
		l.front.Prev = elem
		l.front = elem
	}
	
	l.len++
	return elem
}

func (l *list) PushBack(v interface{}) *ListItem {
	elem := &ListItem{Value: v}

	switch l.Len() {
	case 0:
		l.front = elem
		l.back = elem
	case 1:
		l.back = elem
		l.back.Prev = l.front
		l.front.Next = l.back
	default:
		elem.Prev = l.back
		l.back.Next = elem
		l.back = elem
	}

	l.len++
	return elem
}

func (l *list) Remove(i *ListItem) {
	switch {
	case l.Len() == 1:
		l.front = nil
		l.back = nil
	case i == l.front:
		l.front = l.front.Next
		l.front.Prev = nil
	case i == l.back:
		l.back = l.back.Prev
		l.back.Next = nil
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}
