package db

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/xpmatteo/todomvc-golang/todo"
	"testing"
)

//goland:noinspection SqlNoDataSourceInspection
const createTable = `
create table todo_items (
    id varchar(10),
    title varchar(200),
    isDone bool,
    primary key(id)
);
`

//goland:noinspection SqlNoDataSourceInspection
func Test_readTodoItem_ok(t *testing.T) {
	assert := assert.New(t)
	db := initTestDb()
	mustExec(db, "insert into todo_items (id, title, isDone) values (?, ?, ?)", "123", "foo", false)
	repo := NewTodoRepository(db)
	id := todo.MustNewItemId("123")

	actual, ok, err := repo.Find(id)

	assert.NoError(err)
	assert.True(ok, "got OK from Find")
	expected := &todo.Item{
		Title:  "foo",
		IsDone: false,
		Id:     id,
	}
	assert.Equal(expected, actual)
}

func Test_readTodoItem_notFound(t *testing.T) {
	assert := assert.New(t)
	db := initTestDb()
	repo := NewTodoRepository(db)
	id := todo.MustNewItemId("678")

	item, ok, err := repo.Find(id)

	assert.NoError(err)
	assert.False(ok, "got not OK from Find")
	assert.Nil(item, "got nil for an *item")
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
