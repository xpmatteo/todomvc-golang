package db

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/xpmatteo/todomvc-golang/todo"
	"testing"
)

func testDb() *sql.DB {
	db, err := sql.Open("sqlite", "test.db")
	if err != nil {
		panic(err.Error())
	}
	return db
}

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
func Test_readTodoItem(t *testing.T) {
	assert := assert.New(t)
	db := testDb()
	mustExec(db, "drop table if exists todo_items")
	mustExec(db, createTable)
	mustExec(db, "insert into todo_items (id, title, isDone) values (?, ?, ?)", "123", "foo", false)
	repo := NewTodoRepository(db)
	id := todo.MustNewItemId("123")

	actual, err := repo.Find(id)

	assert.NoError(err)
	expected := &todo.Item{
		Title:  "foo",
		IsDone: false,
		Id:     id,
	}
	assert.Equal(expected, actual)
}

func mustExec(db *sql.DB, sql string, args ...any) {
	_, err := db.Exec(sql, args...)
	if err != nil {
		panic(err.Error())
	}
}
