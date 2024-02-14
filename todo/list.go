package todo

import (
	"errors"
	"fmt"
)

var (
	UserError  = errors.New("user error")
	ErrorBadId = fmt.Errorf("%w: bad todoItemId", UserError)
)

type Item struct {
	Title       string
	IsCompleted bool
	Id          ItemId
	IsDestroyed bool // help with persistence
	IsModified  bool // help with persistence
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
		Title:       title,
		IsCompleted: true,
	})
}

func (l *List) Toggle(id ItemId) error {
	item, ok := l.find(id)
	if !ok {
		return ErrorBadId
	}
	item.IsCompleted = !item.IsCompleted
	item.IsModified = true
	return nil
}

func (l *List) Edit(id ItemId, title string) error {
	item, ok := l.find(id)
	if !ok || item.IsDestroyed {
		return ErrorBadId
	}
	if len(title) == 0 {
		item.IsDestroyed = true
	} else {
		item.Title = title
		item.IsModified = true
	}
	return nil
}

func (l *List) Destroy(id ItemId) error {
	item, ok := l.find(id)
	if ok {
		item.IsDestroyed = true
	}
	return nil
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
		if item.IsCompleted {
			result = append(result, item)
		}
	})
	return result
}

func (l *List) ActiveItems() []*Item {
	var result []*Item
	l.forEach(func(item *Item) {
		if !item.IsCompleted {
			result = append(result, item)
		}
	})
	return result
}

func (l *List) forEach(f func(*Item)) {
	for _, item := range l.Items {
		if item.IsDestroyed {
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
