package web

import (
	"github.com/xpmatteo/todomvc-golang/todo"
	"html/template"
	"net/http"
)

func MakeIndexHandler(templ *template.Template, model *todo.List) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "/active" && r.URL.Path != "/completed" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		data := makeDataForTemplate(model, r)
		executeTemplate(w, templ, data)
	})
}

func MakeNewItemHandler(templ *template.Template, model *todo.List) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			badRequest(w, err)
			return
		}
		model.Add(r.Form.Get("new-todo"))
		data := makeDataForTemplate(model, r)
		executeTemplate(w, templ, data)
	})
}

func MakeToggleHandler(templ *template.Template, model *todo.List) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

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
		err = model.Toggle(id)
		if err != nil {
			badRequest(w, err)
			return
		}

		data := makeDataForTemplate(model, r)
		executeTemplate(w, templ, data)
	})
}

func MakeEditHandler(templ *template.Template, model *todo.List) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
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
		title := r.Form.Get("todoItemTitle")
		if len(title) == 0 {
			model.Destroy(id)
		} else {
			err = model.Edit(id, title)
			if err != nil {
				badRequest(w, err)
				return
			}
		}

		data := makeDataForTemplate(model, r)
		executeTemplate(w, templ, data)
	})
}

func MakeDestroyHandler(templ *template.Template, model *todo.List) http.Handler {
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
		model.Destroy(id)

		data := makeDataForTemplate(model, r)
		executeTemplate(w, templ, data)
	})
}

func badRequest(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}

func makeDataForTemplate(model *todo.List, r *http.Request) map[string]interface{} {
	items := model.AllItems()
	if r.URL.Path == "/completed" {
		items = model.CompletedItems()
	} else if r.URL.Path == "/active" {
		items = model.ActiveItems()
	}
	return map[string]interface{}{
		"Items":         items,
		"Model":         model,
		"Path":          r.URL.Path,
		"EditingItemId": r.URL.Query().Get("edit"),
	}
}

func executeTemplate(w http.ResponseWriter, templ *template.Template, data map[string]interface{}) {
	err := templ.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
