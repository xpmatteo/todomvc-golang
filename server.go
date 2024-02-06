package main

import (
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
	http.Handle("/", web.Logging(web.MakeIndexHandler(templ, model)))
	http.Handle("/new-todo", web.Slowdown(1000, web.Logging(web.MakeNewItemHandler(templ, model))))
	http.Handle("/toggle", web.Logging(web.MakeToggleHandler(templ, model)))
	http.Handle("/edit", web.Logging(web.MakeEditHandler(templ, model)))
	http.Handle("/destroy", web.Logging(web.MakeDestroyHandler(templ, model)))

	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./public/img"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./public/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./public/js"))))

	log.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
