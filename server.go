package main

import (
	"database/sql"
	"github.com/dlmiddlecote/sqlstats"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/xpmatteo/todomvc-golang/db"
	"github.com/xpmatteo/todomvc-golang/web"
	"html/template"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"time"
)

const port = "8080"

func main() {
	pool, err := sql.Open("sqlite", "development.db")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	pool.SetConnMaxLifetime(60 * time.Minute)
	pool.SetMaxIdleConns(3)
	pool.SetMaxOpenConns(10)
	if _, err := pool.Exec(db.CreateTableSQL); err != nil {
		panic(err.Error())
	}
	repository := db.NewTodoRepository(pool)

	// publish DB stats with Prometheus
	collector := sqlstats.NewStatsCollector("todo_db", pool)
	prometheus.MustRegister(collector)

	// I have to use a new ServeMux because the default mux is polluted by
	// expvar, a package that I have no idea how it gets included in the project.
	// Problem is that expvar declares a route for "/debug/vars" that
	// conflicts with the route "/", pending a bugfix in go 1.22.*
	mux := http.NewServeMux()
	templ := template.Must(template.ParseFiles("templates/index.html"))
	mux.Handle("GET /{$}",
		web.Metrics("index",
			web.Logging(
				web.IndexHandler(templ, repository))))
	mux.Handle("GET /active",
		web.Metrics("active",
			web.Logging(
				web.IndexHandler(templ, repository))))
	mux.Handle("GET /completed",
		web.Metrics("completed",
			web.Logging(
				web.IndexHandler(templ, repository))))
	mux.Handle("POST /new-todo",
		web.Metrics("new-todo",
			web.Logging(
				web.Slowdown(1000,
					web.NewItemHandler(templ, repository)))))
	mux.Handle("POST /toggle",
		web.Metrics("toggle",
			web.Logging(
				web.ToggleHandler(templ, repository))))
	mux.Handle("POST /edit",
		web.Metrics("edit",
			web.Logging(
				web.EditHandler(templ, repository))))
	mux.Handle("POST /destroy",
		web.Metrics("destroy",
			web.Logging(
				web.DestroyHandler(templ, repository))))

	mux.Handle("GET /metrics", promhttp.Handler())
	mux.Handle("GET /", http.FileServer(http.Dir("./public/")))

	log.Println("Listening on port " + port)
	web.GracefulListenAndServe(":"+port, mux)
}
