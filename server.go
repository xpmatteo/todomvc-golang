package main

import (
	"html/template"
	"log"
	"net/http"
	"xpug.it/todo"
	"xpug.it/web"
)

const port = "8080"

var model = todo.NewList()

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
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	err := r.ParseForm()
	if err != nil {
		return400(w, err)
		return
	}
	log.Printf("%s %s %s", r.Method, r.URL, r.Form)

	err = model.Toggle(todo.ItemId(r.Form.Get("todoItemId")))
	if err != nil {
		return400(w, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	err := r.ParseForm()
	if err != nil {
		return400(w, err)
		return
	}
	log.Printf("%s %s %s", r.Method, r.URL, r.Form)

	itemId := todo.ItemId(r.Form.Get("todoItemId"))
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
	model.Add("foo")
	model.Add("bar")

	templ := template.Must(template.ParseFiles("templates/index.html"))
	http.Handle("/", web.Logging(web.MakeIndexHandler(templ, model)))
	http.HandleFunc("/new-todo", newItemHandler)
	http.HandleFunc("/toggle", toggleHandler)
	http.HandleFunc("/edit", editHandler)
	http.Handle("/destroy", web.MakeDestroyHandler(model))

	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./public/img"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./public/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./public/js"))))

	log.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
