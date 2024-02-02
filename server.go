package main

import (
	"html/template"
	"log"
	"net/http"
	"xpug.it/todo"
)

const port = "8080"

var model = todo.NewList()

var templates = template.Must(template.ParseFiles("templates/index.html"))

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL)
	err := templates.ExecuteTemplate(w, "index.html", todos)
	check(err)
}

func main() {
	http.HandleFunc("/", indexHandler)

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./public/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./public/js"))))

	log.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
