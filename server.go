package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"xpug.it/todo"
)

const port = "8080"

var model = todo.NewList()

var templates = template.Must(template.ParseFiles("templates/index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL)
	err := templates.ExecuteTemplate(w, "index.html", &model)
	if err != nil {
		panic(err)
	}
}

func newItemHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return400(w, err)
		return
	}
	log.Printf("%s %s %s", r.Method, r.URL, r.Form)
	model.Add(r.Form.Get("new-todo"))
	http.Redirect(w, r, "/", http.StatusFound)
}

func toggleHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return400(w, err)
		return
	}
	log.Printf("%s %s %s", r.Method, r.URL, r.Form)

	itemId, err := strconv.Atoi(r.Form.Get("todoItemId"))
	if err != nil {
		return400(w, errors.New("not a number"))
		return
	}
	err = model.Toggle(itemId)
	if err != nil {
		return400(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return400(w, err)
		return
	}
	log.Printf("%s %s %s", r.Method, r.URL, r.Form)

	itemId, err := strconv.Atoi(r.Form.Get("todoItemId"))
	if err != nil {
		return400(w, errors.New("not a number"))
		return
	}
	title := r.Form.Get("todoItemTitle")
	err = model.Edit(itemId, title)
	if err != nil {
		return400(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func return400(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/new-todo", newItemHandler)
	http.HandleFunc("/toggle", toggleHandler)
	http.HandleFunc("/edit", editHandler)

	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./public/img"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./public/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./public/js"))))

	log.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
