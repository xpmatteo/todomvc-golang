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
	_, _ = fr.Save(todo.Item{Title: title, IsDone: true})
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
		if item.IsDeleted {
			continue
		}
		itemCopy := *item
		if item.Id == nil {
			newId := todo.MustNewItemId(strconv.Itoa(fr.nextId))
			itemCopy.Id = newId
			fr.nextId++
		}
		itemCopy.IsModified = false
		newItems = append(newItems, &itemCopy)
	}
	fr.Items = newItems
	return nil
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
