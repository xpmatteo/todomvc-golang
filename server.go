package main

import (
	"html/template"
	"log"
	"net/http"
	"todo"
	"web"
)

const port = "8080"

var model = todo.NewList()

func main() {
	model.Add("foo")
	model.Add("bar")

	templ := template.Must(template.ParseFiles("templates/index.html"))
	http.Handle("/", web.Logging(web.MakeIndexHandler(templ, model)))
	http.Handle("/new-todo", web.Slowdown(1000, web.Logging(web.MakeNewItemHandler(model))))
	http.Handle("/toggle", web.Logging(web.MakeToggleHandler(model)))
	http.Handle("/edit", web.Logging(web.MakeEditHandler(model)))
	http.Handle("/destroy", web.Logging(web.MakeDestroyHandler(model)))

	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./public/img"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./public/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./public/js"))))

	log.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
