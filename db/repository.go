package db

import (
	"database/sql"
	"github.com/xpmatteo/todomvc-golang/todo"
	_ "modernc.org/sqlite"
	"strconv"
)

//goland:noinspection SqlNoDataSourceInspection
const CreateTableSQL = `
create table if not exists todo_items (
    id INTEGER PRIMARY KEY,
    title varchar(200),
    isDone bool
);
`

type TodoRepository interface {
	Save(item todo.Item) (todo.ItemId, error)
	FindList() (*todo.List, error)
	Destroy(id todo.ItemId) error
}

type todoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) TodoRepository {
	return todoRepository{db}
}

//goland:noinspection SqlNoDataSourceInspection
func (t todoRepository) FindList() (*todo.List, error) {
	selectSql := `
select title, isDone, id
from todo_items
order by id`

	rows, err := t.db.Query(selectSql)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	result := todo.NewList()
	for rows.Next() {
		var title string
		var isDone bool
		var idInt int
		err := rows.Scan(&title, &isDone, &idInt)
		if err != nil {
			return nil, err
		}
		id, err := todo.NewItemId(strconv.Itoa(idInt))
		if err != nil {
			return nil, err
		}
		item := &todo.Item{
			Title:  title,
			IsDone: isDone,
			Id:     id,
		}
		result.Add1(item)
	}

	return result, nil
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

//goland:noinspection SqlNoDataSourceInspection
func (t todoRepository) Destroy(id todo.ItemId) error {
	destroySql := `delete from todo_items where id = ?`
	_, err := t.db.Exec(destroySql, id)
	return err
}
