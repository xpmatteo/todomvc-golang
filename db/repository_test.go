package db

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xpmatteo/todomvc-golang/todo"
	"testing"
)

//goland:noinspection SqlNoDataSourceInspection
const createTable = `
create table if not exists todo_items (
    id INTEGER PRIMARY KEY,
    title varchar(200),
    isDone bool
);
`

func Test_saveAndFind(t *testing.T) {
	assert := assert.New(t)
	db := initTestDb()
	repo := NewTodoRepository(db)
	original := todo.Item{Title: "hello", IsDone: true}

	newId, err := repo.Insert(original)
	require.NoError(t, err)

	actual, err := repo.FindList()
	require.NoError(t, err)

	foundItems := actual.AllItems()
	assert.Equal(1, len(foundItems))
	assert.Equal(original.Title, foundItems[0].Title)
	assert.Equal(original.IsDone, foundItems[0].IsDone)
	assert.Equal(newId, foundItems[0].Id)
}

func Test_findAll(t *testing.T) {
	assert := assert.New(t)
	db := initTestDb()
	repo := NewTodoRepository(db)
	id0, err := repo.Insert(todo.Item{Title: "first", IsDone: false})
	require.NoError(t, err)
	id1, err := repo.Insert(todo.Item{Title: "second", IsDone: true})
	require.NoError(t, err)

	actual, err := repo.FindList()

	all := actual.AllItems()
	assert.Equal(2, len(all))
	assert.Equal("first", all[0].Title)
	assert.Equal("second", all[1].Title)
	assert.Equal(false, all[0].IsDone)
	assert.Equal(true, all[1].IsDone)
	assert.Equal(id0, all[0].Id)
	assert.Equal(id1, all[1].Id)
}

func Test_destroy_ok(t *testing.T) {
	assert := assert.New(t)
	db := initTestDb()
	repo := NewTodoRepository(db)
	_, err := repo.Insert(todo.Item{Title: "first", IsDone: false})
	require.NoError(t, err)
	id1, err := repo.Insert(todo.Item{Title: "second", IsDone: true})
	require.NoError(t, err)

	err = repo.Destroy(id1)
	require.NoError(t, err)

	list, err := repo.FindList()
	require.NoError(t, err)

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

	found, err := repo.FindList()
	require.NoError(t, err)
	foundItems := found.AllItems()
	assert.Equal(2, len(foundItems))
	assert.Equal("first", foundItems[0].Title)
	assert.Equal("second", foundItems[1].Title)
	assert.NotNil(foundItems[0].Id)
	assert.NotNil(foundItems[1].Id)
}

func Test_saveModifiedList_isDone(t *testing.T) {
	assert := assert.New(t)
	db := initTestDb()
	repo := NewTodoRepository(db)
	id, err := repo.Insert(todo.Item{Title: "any"})
	require.NoError(t, err)
	list, err := repo.FindList()
	require.NoError(t, err)
	err = list.Toggle(id)
	require.NoError(t, err)

	err = repo.SaveList(list)
	require.NoError(t, err)

	found, err := repo.FindList()
	require.NoError(t, err)
	foundItems := found.AllItems()
	assert.Equal(1, len(foundItems))
	assert.Equal(id, foundItems[0].Id)
	assert.Equal(true, foundItems[0].IsDone)
}

func Test_saveModifiedList_editTitle(t *testing.T) {
	assert := assert.New(t)
	db := initTestDb()
	repo := NewTodoRepository(db)
	id, err := repo.Insert(todo.Item{Title: "any"})
	require.NoError(t, err)
	list, err := repo.FindList()
	require.NoError(t, err)

	err = list.Edit(id, "newTitle")
	require.NoError(t, err)
	err = repo.SaveList(list)
	require.NoError(t, err)

	found, err := repo.FindList()
	require.NoError(t, err)
	foundItems := found.AllItems()
	assert.Equal(1, len(foundItems))
	assert.Equal(id, foundItems[0].Id)
	assert.Equal("newTitle", foundItems[0].Title)
}

func Test_saveModifiedList_destroyItem(t *testing.T) {
	assert := assert.New(t)
	db := initTestDb()
	repo := NewTodoRepository(db)
	id, err := repo.Insert(todo.Item{Title: "any"})
	require.NoError(t, err)
	list, err := repo.FindList()
	require.NoError(t, err)

	list.Destroy(id)
	err = repo.SaveList(list)
	require.NoError(t, err)

	found, err := repo.FindList()
	require.NoError(t, err)
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
	mustExec(db, createTable)
	return db
}

func mustExec(db *sql.DB, sql string, args ...any) {
	_, err := db.Exec(sql, args...)
	if err != nil {
		panic(err.Error())
	}
}
