package doublylinkedlist

type Item struct {
	value interface{}
	next  *Item
	prev  *Item
}

func (i Item) Value() interface{} {
	return i.value
}

func (i Item) Next() *Item {
	return i.next
}

func (i Item) Prev() *Item {
	return i.prev
}

type List struct {
	first  *Item // pointer to first node of the list
	last   *Item // pointer to last node of the list
	length int   // the length of the list
}

func (l List) Len() int {
	return l.length
}

func (l List) First() *Item {
	return l.first
}

func (l List) Last() *Item {
	return l.last
}

func (l *List) PushFront(v interface{}) {
	i := &Item{
		value: v,
		next:  l.First(),
	}
	if l.First() != nil {
		l.First().prev = i
	} else {
		// First item in list
		l.last = i
	}
	l.first = i
	l.length++
}

func (l *List) PushBack(v interface{}) {
	i := &Item{
		value: v,
		prev:  l.Last(),
	}
	if l.Last() != nil {
		l.Last().next = i
	} else {
		// First item in list
		l.first = i
	}
	l.last = i
	l.length++
}

func (l *List) Remove(i Item) {
	// Check if item belongs to list with O(1) complexity
	if i.Prev() == nil && i.Next() == nil {
		return
	}

	if i.Prev() == nil {
		l.first = i.Next()
	} else {
		i.Prev().next = i.Next()
	}
	if i.Next() == nil {
		l.last = i.Prev()
	} else {
		i.Next().prev = i.Prev()
	}
	l.length--
}
