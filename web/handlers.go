package web

import (
	"html/template"
	"log"
	"net/http"
	"todo"
)

func MakeIndexHandler(templ *template.Template, model interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		data := map[string]interface{}{
			"Model":         model,
			"EditingItemId": r.URL.Query().Get("edit"),
		}
		err := templ.Execute(w, data)
		if err != nil {
			panic(err.Error())
		}
	})
}

func MakeNewItemHandler(list *todo.List) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			badRequest(w, err)
			return
		}
		list.Add(r.Form.Get("new-todo"))
		http.Redirect(w, r, r.URL.Path, http.StatusFound)
	})
}

func MakeToggleHandler(list *todo.List) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, r.URL.Path, http.StatusTemporaryRedirect)
			return
		}

		err := r.ParseForm()
		if err != nil {
			badRequest(w, err)
			return
		}
		log.Printf("%s %s %s", r.Method, r.URL, r.Form)

		id, err := todo.NewItemId(r.Form.Get("todoItemId"))
		if err != nil {
			badRequest(w, err)
			return
		}
		err = list.Toggle(id)
		if err != nil {
			badRequest(w, err)
			return
		}
		http.Redirect(w, r, r.URL.Path, http.StatusFound)
	})
}

func MakeEditHandler(list *todo.List) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, r.URL.Path, http.StatusTemporaryRedirect)
			return
		}
		err := r.ParseForm()
		if err != nil {
			badRequest(w, err)
			return
		}
		log.Printf("%s %s %s", r.Method, r.URL, r.Form)

		id, err := todo.NewItemId(r.Form.Get("todoItemId"))
		if err != nil {
			badRequest(w, err)
			return
		}
		title := r.Form.Get("todoItemTitle")
		if len(title) == 0 {
			list.Destroy(id)
		} else {
			err = list.Edit(id, title)
			if err != nil {
				badRequest(w, err)
				return
			}
		}

		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
	})
}

func MakeDestroyHandler(list *todo.List) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			badRequest(w, err)
			return
		}
		id, err := todo.NewItemId(r.Form.Get("todoItemId"))
		if err != nil {
			badRequest(w, err)
			return
		}
		list.Destroy(id)
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
	})
}

func badRequest(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}
