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
    isCompleted bool
);
`

type TodoRepository interface {
	Insert(item todo.Item) (todo.ItemId, error)
	FindList() (*todo.List, error)
	SaveList(list *todo.List) error
}

type todoRepository struct {
	db *sql.DB
}

func (t todoRepository) SaveList(list *todo.List) error {
	for _, item := range list.Items {
		if item.IsDestroyed {
			err := t.Destroy(item.Id)
			if err != nil {
				return err
			}
		} else if item.IsModified {
			err := t.Update(item)
			if err != nil {
				return err
			}
		} else if item.Id == nil {
			newId, err := t.Insert(*item)
			// intentionally modifying the passed-in list, so that it can be shown
			// correctly on the ui
			item.Id = newId
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewTodoRepository(db *sql.DB) TodoRepository {
	return todoRepository{db}
}

//goland:noinspection SqlNoDataSourceInspection
func (t todoRepository) FindList() (*todo.List, error) {
	selectSql := `
select title, isCompleted, id
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
		var isCompleted bool
		var idInt int
		err := rows.Scan(&title, &isCompleted, &idInt)
		if err != nil {
			return nil, err
		}
		id, err := todo.NewItemId(strconv.Itoa(idInt))
		if err != nil {
			return nil, err
		}
		item := &todo.Item{
			Title:       title,
			IsCompleted: isCompleted,
			Id:          id,
		}
		result.Add1(item)
	}

	return result, nil
}

//goland:noinspection SqlNoDataSourceInspection
func (t todoRepository) Insert(item todo.Item) (todo.ItemId, error) {
	sql := `
insert into todo_items
	(title, isCompleted)
values (?, ?)`
	//tx, err := t.db.Begin()
	//if err != nil {
	//	return nil, err
	//}
	result, err := t.db.Exec(sql, item.Title, item.IsCompleted)
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

//goland:noinspection SqlNoDataSourceInspection
func (t todoRepository) Update(item *todo.Item) error {
	updateSql := `
update todo_items
set title = ?
  , isCompleted = ?
where id = ?`
	_, err := t.db.Exec(updateSql, item.Title, item.IsCompleted, item.Id)
	return err
}
