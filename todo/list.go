package todo

import (
	"errors"
	"strconv"
)

type Item struct {
	Title  string
	IsDone bool
	Id     ItemId
}

func (item Item) String() string {
	return "Item"
}

type List struct {
	Items  map[ItemId]*Item
	nextId int
}

func NewList() *List {
	return &List{make(map[ItemId]*Item), 0}
}

func (l *List) Add(title string) ItemId {
	if len(title) == 0 {
		return nil
	}
	newId := MustNewItemId(strconv.Itoa(l.nextId))
	l.Items[newId] = &Item{title, false, newId}
	l.nextId++
	return newId
}

func (l *List) AddCompleted(title string) ItemId {
	newId := l.Add(title)
	if newId != nil {
		_ = l.Toggle(newId)
	}
	return newId
}

func (l *List) Toggle(id ItemId) error {
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

func (l *List) Edit(id ItemId, title string) error {
	item, ok := l.Items[id]
	if !ok {
		return errors.New("bad todo-item ID")
	}
	item.Title = title
	return nil
}

func (l *List) Destroy(id ItemId) {
	delete(l.Items, id)
}

func (l *List) AllItems() []*Item {
	result := []*Item{}
	for _, item := range l.Items {
		result = append(result, item)
	}
	return result
}

func (l *List) CompletedItems() []*Item {
	result := []*Item{}
	for _, item := range l.Items {
		if item.IsDone {
			result = append(result, item)
		}
	}
	return result
}

func (l *List) ActiveItems() []*Item {
	result := []*Item{}
	for _, item := range l.Items {
		if !item.IsDone {
			result = append(result, item)
		}
	}
	return result
}
