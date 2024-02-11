package db

import (
	"database/sql"
	"errors"
	"github.com/xpmatteo/todomvc-golang/todo"
	_ "modernc.org/sqlite"
)

type TodoRepository interface {
	Find(todo.ItemId) (*todo.Item, error)
}

type todoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) TodoRepository {
	return todoRepository{db}
}

//goland:noinspection SqlNoDataSourceInspection
func (t todoRepository) Find(id todo.ItemId) (*todo.Item, error) {
	sql := `
select title, isDone
from todo_items
where id = ?`
	rows, err := t.db.Query(sql, id.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var title string
		var isDone bool
		rows.Scan(&title, &isDone)
		return &todo.Item{
			Title:  title,
			IsDone: isDone,
			Id:     id,
		}, nil
	}

	return nil, errors.New("not found")
}
