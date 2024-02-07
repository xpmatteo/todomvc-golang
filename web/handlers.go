package web

import (
	"fmt"
	"github.com/xpmatteo/todomvc-golang/todo"
	"html/template"
	"net/http"
)

const (
	pathActive       = "/active"
	pathCompleted    = "/completed"
	keyTodoItemId    = "todoItemId"
	keyTodoItemTitle = "todoItemTitle"
)

func MakeIndexHandler(templ *template.Template, model *todo.List) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != pathActive && r.URL.Path != pathCompleted {
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

		id, err := todo.NewItemId(r.Form.Get(keyTodoItemId))
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

		id, err := todo.NewItemId(r.Form.Get(keyTodoItemId))
		if err != nil {
			badRequest(w, err)
			return
		}
		title := r.Form.Get(keyTodoItemTitle)
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
		id, err := todo.NewItemId(r.Form.Get(keyTodoItemId))
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
	path := determinePath(r)
	if path == pathCompleted {
		items = model.CompletedItems()
	} else if path == pathActive {
		items = model.ActiveItems()
	}
	return map[string]interface{}{
		"Items":            items,
		"Path":             path,
		"ItemsCount":       len(model.Items),
		"NoCompletedItems": len(model.CompletedItems()) == 0,
		"ItemsLeft":        itemsLeft(model),
		"EditingItemId":    r.URL.Query().Get("edit"),
	}
}

func itemsLeft(model *todo.List) string {
	count := len(model.ActiveItems())
	if count == 1 {
		return "1 item left"
	}
	return fmt.Sprintf("%d items left", count)
}

func determinePath(r *http.Request) string {
	path := r.URL.Path
	_ = r.ParseForm()
	pathFromForm := r.Form.Get("Path")
	if pathFromForm != "" {
		path = pathFromForm
	}
	return path
}

func executeTemplate(w http.ResponseWriter, templ *template.Template, data map[string]interface{}) {
	err := templ.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
