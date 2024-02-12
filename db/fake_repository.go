package db

import (
	"github.com/xpmatteo/todomvc-golang/todo"
	"strconv"
)

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

type FakeRepositoryImplementation struct {
	items  []*todo.Item
	nextId int
}

func (fr *FakeRepositoryImplementation) Save(item todo.Item) (todo.ItemId, error) {
	newId := todo.MustNewItemId(strconv.Itoa(fr.nextId))
	item.Id = newId
	fr.items = append(fr.items, &item)
	fr.nextId++
	return newId, nil
}

func (fr *FakeRepositoryImplementation) FindList() (*todo.List, error) {
	result := todo.NewList()
	for _, item := range fr.items {
		result.Add1(item)
	}
	return result, nil
}

func (fr *FakeRepositoryImplementation) Destroy(id todo.ItemId) error {
	var newItems []*todo.Item
	for _, item := range fr.items {
		if item.Id != id {
			newItems = append(newItems, item)
		}
	}
	fr.items = newItems
	return nil
}
