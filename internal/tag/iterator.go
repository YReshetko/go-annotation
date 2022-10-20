package tag

import "reflect"

type iterator struct {
	t reflect.Type

	peekIndex int
	nextIndex int
	isHasNext bool
}

func newIterator(t reflect.Type) *iterator {
	i := &iterator{
		t:         t,
		peekIndex: -1,
	}
	i.findNext()
	return i
}

func (i *iterator) hasNext() bool {
	return i.isHasNext
}

func (i *iterator) next() (reflect.StructField, int) {
	if i.peekIndex != i.nextIndex {
		i.peekIndex = i.nextIndex
	}
	i.findNext()
	return i.t.Field(i.peekIndex), i.peekIndex
}

func (i *iterator) findNext() {
	i.nextIndex = -1
	i.isHasNext = false

	n := i.t.NumField()
	for j := i.peekIndex + 1; j < n; j++ {
		if !i.t.Field(j).IsExported() {
			continue
		}
		i.nextIndex = j
		i.isHasNext = true
		return
	}
}
