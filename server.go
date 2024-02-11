package main

import (
	"database/sql"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/xpmatteo/todomvc-golang/db"
	"github.com/xpmatteo/todomvc-golang/todo"
	"github.com/xpmatteo/todomvc-golang/web"
	"html/template"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"time"
)

const port = "8080"

var model = todo.NewList()

//goland:noinspection SqlNoDataSourceInspection
const createTable = `
create table if not exists todo_items (
    id INTEGER PRIMARY KEY,
    title varchar(200),
    isDone bool
);
`

func main() {
	pool, err := sql.Open("sqlite", "development.db")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	pool.SetConnMaxLifetime(60 * time.Minute)
	pool.SetMaxIdleConns(3)
	pool.SetMaxOpenConns(10)
	if _, err := pool.Exec(createTable); err != nil {
		panic(err.Error())
	}
	repository := db.NewTodoRepository(pool)
	InsertTodo(repository, "foo", false)
	InsertTodo(repository, "bar", true)
	InsertTodo(repository, "baz", false)

	templ := template.Must(template.ParseFiles("templates/index.html"))
	http.Handle("/",
		web.Metrics("index",
			web.Logging(
				web.GETonly(
					web.IndexHandler(templ, model)))))
	http.Handle("/new-todo",
		web.Metrics("new-todo",
			web.Logging(
				web.POSTonly(
					web.Slowdown(1000,
						web.NewItemHandler(templ, model))))))
	http.Handle("/toggle",
		web.Metrics("toggle",
			web.Logging(
				web.POSTonly(
					web.ToggleHandler(templ, model)))))
	http.Handle("/edit",
		web.Metrics("edit",
			web.Logging(
				web.POSTonly(
					web.EditHandler(templ, model)))))
	http.Handle("/destroy",
		web.Metrics("destroy",
			web.Logging(
				web.POSTonly(
					web.DestroyHandler(templ, model)))))

	http.Handle("/metrics", promhttp.Handler())

	web.GET("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./public/img"))))
	web.GET("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./public/css"))))
	web.GET("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./public/js"))))

	log.Println("Listening on port " + port)
	web.GracefulListenAndServe(":"+port, nil)
}

func InsertTodo(repository db.TodoRepository, title string, isDone bool) {
	_, err := repository.Save(todo.Item{
		Title:  title,
		IsDone: isDone,
	})
	if err != nil {
		panic(err.Error())
	}
}
