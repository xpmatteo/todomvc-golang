package db

import (
	"database/sql"
	"github.com/xpmatteo/todomvc-golang/todo"
	_ "modernc.org/sqlite"
	"strconv"
)

type TodoRepository interface {
	Find(todo.ItemId) (*todo.Item, bool, error)
	Save(item todo.Item) (todo.ItemId, error)
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
	rows, err := t.db.Query(sql, id)
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

//goland:noinspection SqlNoDataSourceInspection
func (t todoRepository) Save(item todo.Item) (todo.ItemId, error) {
	sql := `
insert into todo_items
	(title, isDone)
values (?, ?)`
	//tx, err := t.db.Begin()
	//if err != nil {
	//	return nil, err
	//}
	result, err := t.db.Exec(sql, item.Title, item.IsDone)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	//if err := tx.Commit(); err != nil {
	//	return nil, err
	//}

	newId, err := todo.NewItemId(strconv.FormatInt(id, 10))
	if err != nil {
		return nil, err
	}
	return newId, nil
}
