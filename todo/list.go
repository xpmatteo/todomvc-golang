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

type List struct {
	Items  map[ItemId]*Item
	ids    []ItemId
	nextId int
}

func NewList() *List {
	return &List{
		make(map[ItemId]*Item),
		make([]ItemId, 0),
		0}
}

func (l *List) Add(title string) ItemId {
	if len(title) == 0 {
		return nil
	}
	newId := MustNewItemId(strconv.Itoa(l.nextId))
	l.Items[newId] = &Item{title, false, newId}
	l.nextId++
	l.ids = append(l.ids, newId)
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

func (l *List) Edit(id ItemId, title string) error {
	item, ok := l.Items[id]
	if !ok {
		return errors.New("bad todo-item ID")
	}
	item.Title = title
	return nil
}

func (l *List) Destroy(id ItemId) {
	// we do not update l.ids; we will simply ignore missing items in forEach
	delete(l.Items, id)
}

func (l *List) AllItems() []*Item {
	result := []*Item{}
	l.forEach(func(item *Item) {
		result = append(result, item)
	})
	return result
}

func (l *List) CompletedItems() []*Item {
	result := []*Item{}
	l.forEach(func(item *Item) {
		if item.IsDone {
			result = append(result, item)
		}
	})
	return result
}

func (l *List) ActiveItems() []*Item {
	result := []*Item{}
	l.forEach(func(item *Item) {
		if !item.IsDone {
			result = append(result, item)
		}
	})
	return result
}

func (l *List) forEach(f func(*Item)) {
	for _, id := range l.ids {
		item := l.Items[id]
		if item == nil {
			// item was destroyed
			continue
		}
		f(item)
	}
}
