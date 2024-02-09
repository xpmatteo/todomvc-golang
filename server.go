package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/xpmatteo/todomvc-golang/todo"
	"github.com/xpmatteo/todomvc-golang/web"
	"html/template"
	"log"
	"net/http"
)

const port = "8080"

var model = todo.NewList()

func main() {
	model.Add("foo")
	model.AddCompleted("bar")
	model.Add("baz")

	templ := template.Must(template.ParseFiles("templates/index.html"))
	web.GET("/",
		web.Metrics("index",
			web.Logging(web.IndexHandler(templ, model))))
	web.POST("/new-todo",
		web.Metrics("new-todo",
			web.Slowdown(1000,
				web.Logging(web.NewItemHandler(templ, model)))))
	web.POST("/toggle",
		web.Metrics("toggle",
			web.Logging(
				web.ToggleHandler(templ, model))))
	web.POST("/edit",
		web.Logging(web.EditHandler(templ, model)))
	web.POST("/destroy",
		web.Logging(web.DestroyHandler(templ, model)))

	http.Handle("/metrics", promhttp.Handler())

	web.GET("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./public/img"))))
	web.GET("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./public/css"))))
	web.GET("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./public/js"))))

	log.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
