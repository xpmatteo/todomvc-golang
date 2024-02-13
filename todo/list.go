package todo

import (
	"errors"
)

type Item struct {
	Title      string
	IsDone     bool
	Id         ItemId
	IsDeleted  bool // help with persistence
	IsModified bool // help with persistence
}

type List struct {
	Items  []*Item
	nextId int
}

func NewList() *List {
	return &List{
		make([]*Item, 0),
		0}
}

func (l *List) Add(title string, id ItemId) {
	if len(title) == 0 {
		return
	}
	l.Items = append(l.Items, &Item{Title: title, Id: id})
}

func (l *List) Add1(item *Item) {
	l.Items = append(l.Items, item)
}

func (l *List) AddCompleted(title string) {
	l.Add1(&Item{
		Title:  title,
		IsDone: true,
	})
}

func (l *List) Toggle(id ItemId) error {
	item, ok := l.find(id)
	if !ok {
		return errors.New("bad todo-item ID")
	}
	item.IsDone = !item.IsDone
	item.IsModified = true
	return nil
}

func (l *List) Edit(id ItemId, title string) error {
	item, ok := l.find(id)
	if !ok || item.IsDeleted {
		return errors.New("bad todo-item ID")
	}
	if len(title) == 0 {
		item.IsDeleted = true
	} else {
		item.Title = title
		item.IsModified = true
	}
	return nil
}

func (l *List) Destroy(id ItemId) {
	item, ok := l.find(id)
	if ok {
		item.IsDeleted = true
	}
}

func (l *List) AllItems() []*Item {
	var result []*Item
	l.forEach(func(item *Item) {
		result = append(result, item)
	})
	return result
}

func (l *List) CompletedItems() []*Item {
	var result []*Item
	l.forEach(func(item *Item) {
		if item.IsDone {
			result = append(result, item)
		}
	})
	return result
}

func (l *List) ActiveItems() []*Item {
	var result []*Item
	l.forEach(func(item *Item) {
		if !item.IsDone {
			result = append(result, item)
		}
	})
	return result
}

func (l *List) forEach(f func(*Item)) {
	for _, item := range l.Items {
		if item.IsDeleted {
			continue
		}
		f(item)
	}
}

func (l *List) find(id ItemId) (item *Item, ok bool) {
	for _, item := range l.Items {
		if item.Id == id {
			return item, true
		}
	}
	return nil, false
}
