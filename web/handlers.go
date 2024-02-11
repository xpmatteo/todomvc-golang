package web

import (
	"github.com/xpmatteo/todomvc-golang/db"
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

func IndexHandler(templ *template.Template, repo db.TodoRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != pathActive && r.URL.Path != pathCompleted {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		model, err := repo.FindList()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		print(model.Items)
		vm := viewModel(model, r)
		render(w, r, templ, vm)
	})
}

func NewItemHandler(templ *template.Template, model *todo.List) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			badRequest(w, err)
			return
		}
		model.Add(r.Form.Get("new-todo"))
		vm := viewModel(model, r)
		render(w, r, templ, vm)
	})
}

func ToggleHandler(templ *template.Template, model *todo.List) http.Handler {
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
		err = model.Toggle(id)
		if err != nil {
			badRequest(w, err)
			return
		}

		vm := viewModel(model, r)
		render(w, r, templ, vm)
	})
}

func EditHandler(templ *template.Template, model *todo.List) http.Handler {
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

		vm := viewModel(model, r)
		render(w, r, templ, vm)
	})
}

func DestroyHandler(templ *template.Template, model *todo.List) http.Handler {
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

		vm := viewModel(model, r)
		render(w, r, templ, vm)
	})
}

func badRequest(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}
