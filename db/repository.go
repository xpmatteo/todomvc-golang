package db

import (
	"database/sql"
	"github.com/xpmatteo/todomvc-golang/todo"
	_ "modernc.org/sqlite"
)

type TodoRepository interface {
	Find(todo.ItemId) (*todo.Item, bool, error)
}

type todoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) TodoRepository {
	return todoRepository{db}
}

//goland:noinspection SqlNoDataSourceInspection
func (t todoRepository) Find(id todo.ItemId) (item *todo.Item, ok bool, err error) {
	sql := `
select title, isDone
from todo_items
where id = ?`
	rows, err := t.db.Query(sql, id.String())
	if err != nil {
		return nil, false, err
	}
	defer func() { _ = rows.Close() }()

	if rows.Next() {
		var title string
		var isDone bool
		err := rows.Scan(&title, &isDone)
		if err != nil {
			return nil, false, err
		}
		return &todo.Item{
			Title:  title,
			IsDone: isDone,
			Id:     id,
		}, true, nil
	}

	return nil, false, nil
}
