package db

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xpmatteo/todomvc-golang/todo"
	"testing"
)

func Test_saveAndFind(t *testing.T) {
	assert := assert.New(t)
	db := initTestDb()
	repo := NewTodoRepository(db)
	original := todo.Item{Title: "hello", IsCompleted: true}

	newId := mustInsert(repo, &original)
	actual := mustFindList(repo)

	foundItems := actual.AllItems()
	assert.Equal(1, len(foundItems))
	assert.Equal(original.Title, foundItems[0].Title)
	assert.Equal(original.IsCompleted, foundItems[0].IsCompleted)
	assert.Equal(newId, foundItems[0].Id)
}

func Test_findAll(t *testing.T) {
	assert := assert.New(t)
	db := initTestDb()
	repo := NewTodoRepository(db)
	id0 := mustInsert(repo, &todo.Item{Title: "first", IsCompleted: false})
	id1 := mustInsert(repo, &todo.Item{Title: "second", IsCompleted: true})

	actual := mustFindList(repo)

	all := actual.AllItems()
	assert.Equal(2, len(all))
	assert.Equal("first", all[0].Title)
	assert.Equal("second", all[1].Title)
	assert.Equal(false, all[0].IsCompleted)
	assert.Equal(true, all[1].IsCompleted)
	assert.Equal(id0, all[0].Id)
	assert.Equal(id1, all[1].Id)
}

func Test_destroy_ok(t *testing.T) {
	assert := assert.New(t)
	db := initTestDb()
	repo := todoRepository{db}
	_ = mustInsert(repo, &todo.Item{Title: "first", IsCompleted: false})
	id1 := mustInsert(repo, &todo.Item{Title: "second", IsCompleted: true})

	err := repo.Destroy(id1)
	require.NoError(t, err)

	list := mustFindList(repo)
	assert.Equal(1, len(list.Items))
	assert.Equal("first", list.Items[0].Title)
}

func Test_saveNewList(t *testing.T) {
	assert := assert.New(t)
	db := initTestDb()
	repo := NewTodoRepository(db)
	list := todo.NewList()
	list.Add("first", nil)
	list.Add("second", nil)

	err := repo.SaveList(list)
	require.NoError(t, err)

	found := mustFindList(repo)
	foundItems := found.AllItems()
	assert.Equal(2, len(foundItems))
	assert.Equal("first", foundItems[0].Title)
	assert.Equal("second", foundItems[1].Title)
	assert.NotNil(foundItems[0].Id)
	assert.NotNil(foundItems[1].Id)
	// the passed-in list now has ids populated
	assert.Equal(foundItems[0].Id, list.Items[0].Id)
	assert.Equal(foundItems[1].Id, list.Items[1].Id)
}

func Test_saveModifiedList_isDone(t *testing.T) {
	assert := assert.New(t)
	repo := NewTodoRepository(initTestDb())
	id := mustInsert(repo, &todo.Item{Title: "any"})
	list, err := repo.FindList()
	require.NoError(t, err)
	err = list.Toggle(id)
	require.NoError(t, err)

	err = repo.SaveList(list)
	require.NoError(t, err)

	found := mustFindList(repo)
	foundItems := found.AllItems()
	assert.Equal(1, len(foundItems))
	assert.Equal(id, foundItems[0].Id)
	assert.Equal(true, foundItems[0].IsCompleted)
}

func Test_saveModifiedList_editTitle(t *testing.T) {
	assert := assert.New(t)
	repo := NewTodoRepository(initTestDb())
	id := mustInsert(repo, &todo.Item{Title: "any"})
	list, err := repo.FindList()
	require.NoError(t, err)

	err = list.Edit(id, "newTitle")
	require.NoError(t, err)
	mustSaveList(repo, list)

	found := mustFindList(repo)
	foundItems := found.AllItems()
	assert.Equal(1, len(foundItems))
	assert.Equal(id, foundItems[0].Id)
	assert.Equal("newTitle", foundItems[0].Title)
}

func Test_saveModifiedList_destroyItem(t *testing.T) {
	assert := assert.New(t)
	repo := NewTodoRepository(initTestDb())
	id := mustInsert(repo, &todo.Item{Title: "any"})
	list := mustFindList(repo)

	list.Destroy(id)
	mustSaveList(repo, list)

	found := mustFindList(repo)
	foundItems := found.AllItems()
	assert.Equal(0, len(foundItems))
}

//goland:noinspection SqlNoDataSourceInspection
func initTestDb() *sql.DB {
	db, err := sql.Open("sqlite", "test.db")
	if err != nil {
		panic(err.Error())
	}
	mustExec(db, "drop table if exists todo_items")
	mustExec(db, CreateTableSQL)
	return db
}

func mustExec(db *sql.DB, sql string, args ...any) {
	_, err := db.Exec(sql, args...)
	if err != nil {
		panic(err.Error())
	}
}

func mustInsert(repo TodoRepository, item *todo.Item) todo.ItemId {
	id, err := repo.Insert(*item)
	if err != nil {
		panic(err.Error())
	}
	return id
}

func mustFindList(repo TodoRepository) *todo.List {
	list, err := repo.FindList()
	if err != nil {
		panic(err.Error())
	}
	return list
}

func mustSaveList(repo TodoRepository, list *todo.List) {
	err := repo.SaveList(list)
	if err != nil {
		panic(err.Error())
	}
}
