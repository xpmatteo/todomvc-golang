package todo

import (
	"errors"
	"strconv"
)

type Item struct {
	Title  string
	IsDone bool
	Id     itemId
}

type List struct {
	Items  map[itemId]*Item
	nextId int
}

func NewList() *List {
	return &List{make(map[itemId]*Item), 0}
}

func (l *List) Add(title string) {
	if len(title) == 0 {
		return
	}
	newId := itemId(strconv.Itoa(l.nextId))
	l.Items[newId] = &Item{title, false, newId}
	l.nextId++
}

func (l *List) Toggle(id itemId) error {
	item, ok := l.Items[id]
	if !ok {
		return errors.New("bad todo-item ID")
	}
	item.IsDone = !item.IsDone
	return nil
}

func (l *List) ItemsLeft() int {
	result := 0
	for _, item := range l.Items {
		if !item.IsDone {
			result++
		}
	}
	return result
}

func (l *List) Edit(id itemId, title string) error {
	item, ok := l.Items[id]
	if !ok {
		return errors.New("bad todo-item ID")
	}
	item.Title = title
	return nil
}

func (l *List) Destroy(id itemId) {
	delete(l.Items, id)
}
