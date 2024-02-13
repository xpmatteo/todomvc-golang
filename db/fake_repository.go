package db

import (
	"github.com/xpmatteo/todomvc-golang/todo"
	"strconv"
)

type FakeRepositoryImplementation struct {
	Items  []*todo.Item
	nextId int
}

func FakeRepository() *FakeRepositoryImplementation {
	return &FakeRepositoryImplementation{}
}

func (fr *FakeRepositoryImplementation) Add(title string) *FakeRepositoryImplementation {
	_, _ = fr.Save(todo.Item{Title: title})
	return fr
}

func (fr *FakeRepositoryImplementation) AddCompleted(title string) *FakeRepositoryImplementation {
	_, _ = fr.Save(todo.Item{Title: title, IsCompleted: true})
	return fr
}

func (fr *FakeRepositoryImplementation) Save(item todo.Item) (todo.ItemId, error) {
	newId := todo.MustNewItemId(strconv.Itoa(fr.nextId))
	item.Id = newId
	fr.Items = append(fr.Items, &item)
	fr.nextId++
	return newId, nil
}

func (fr *FakeRepositoryImplementation) FindList() (*todo.List, error) {
	result := todo.NewList()
	for _, item := range fr.Items {
		result.Add1(item)
	}
	return result, nil
}

func (fr *FakeRepositoryImplementation) SaveList(list *todo.List) error {
	var newItems []*todo.Item
	for _, item := range list.Items {
		if item.IsDestroyed {
			continue
		}
		// intentionally modifying the passed in list
		if item.Id == nil {
			item.Id = fr.newId()
		}
		item.IsModified = false
		newItems = append(newItems, item)
	}
	fr.Items = newItems
	return nil
}

func (fr *FakeRepositoryImplementation) newId() todo.ItemId {
	newId := todo.MustNewItemId(strconv.Itoa(fr.nextId))
	fr.nextId++
	return newId
}

func (fr *FakeRepositoryImplementation) Destroy(id todo.ItemId) error {
	var newItems []*todo.Item
	for _, item := range fr.Items {
		if item.Id != id {
			newItems = append(newItems, item)
		}
	}
	fr.Items = newItems
	return nil
}

func (fr *FakeRepositoryImplementation) Insert(item todo.Item) error {
	newItem := item
	newItem.Id = fr.newId()
	fr.Items = append(fr.Items, &newItem)
	return nil
}
