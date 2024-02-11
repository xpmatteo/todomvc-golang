package main

import (
	"context"
	"database/sql"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

func main() {
	db, err := sql.Open("sqlite", "development.db")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetMaxIdleConns(3)
	db.SetMaxOpenConns(10)

	if err := ping(db); err != nil {
		log.Fatal(err.Error())
		return
	}

	model.Add("foo")
	model.AddCompleted("bar")
	model.Add("baz")

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

func ping(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return db.PingContext(ctx)
}
