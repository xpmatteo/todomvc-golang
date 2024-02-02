package main

import (
	"encoding/json"
	"fmt"
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
	err := templates.ExecuteTemplate(w, "index.html", model)
	check(err)
}

func newItemHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Error: %s", err)
		_, _ = fmt.Fprintf(w, "oh no %s", err)
		return
	}
	log.Printf("%s %s %s", r.Method, r.URL, r.Form)
	model.Add(r.Form.Get("new-todo"))
	http.Redirect(w, r, "/", http.StatusFound)
}

func toggleHandler(w http.ResponseWriter, r *http.Request) {
	var itemId int
	err := json.NewDecoder(r.Body).Decode(&itemId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = model.Toggle(itemId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
	log.Printf("%s %s %d", r.Method, r.URL, itemId)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/new-todo", newItemHandler)
	http.HandleFunc("/toggle", toggleHandler)

	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./public/img"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./public/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./public/js"))))

	log.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
