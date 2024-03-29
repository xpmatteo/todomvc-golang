package web

import (
	"errors"
	"github.com/xpmatteo/todomvc-golang/todo"
	"html/template"
	"log"
	"net/http"
)

const (
	pathActive       = "/active"
	pathCompleted    = "/completed"
	keyTodoItemId    = "todoItemId"
	keyTodoItemTitle = "todoItemTitle"
)

type ListFinder interface {
	FindList() (*todo.List, error)
}

type Repository interface {
	FindList() (*todo.List, error)
	SaveList(*todo.List) error
}

func IndexHandler(templ *template.Template, repo ListFinder) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		model, err := repo.FindList()
		if err != nil {
			internalServerError(w, err)
			return
		}
		vm := viewModel(model, r)
		render(w, r, templ, vm)
	})
}

func NewItemHandler(templ *template.Template, repository Repository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			badRequest(w, err)
			return
		}

		title := r.Form.Get("new-todo")
		model, err := with(repository, func(model *todo.List) error {
			model.Add(title, nil)
			return nil
		})
		if err != nil {
			handleError(w, err)
			return
		}

		vm := viewModel(model, r)
		render(w, r, templ, vm)
	})
}

func ToggleHandler(templ *template.Template, repository Repository) http.Handler {
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
		model, err := with(repository, func(model *todo.List) error {
			return model.Toggle(id)
		})
		if err != nil {
			handleError(w, err)
			return
		}
		vm := viewModel(model, r)
		render(w, r, templ, vm)
	})
}

func handleError(w http.ResponseWriter, err error) {
	if errors.Is(err, todo.UserError) {
		badRequest(w, err)
	} else {
		internalServerError(w, err)
	}
}

func EditHandler(templ *template.Template, repository Repository) http.Handler {
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

		model, err := with(repository, func(model *todo.List) error {
			return model.Edit(id, title)
		})
		if err != nil {
			handleError(w, err)
			return
		}

		vm := viewModel(model, r)
		render(w, r, templ, vm)
	})
}

func DestroyHandler(templ *template.Template, repository Repository) http.Handler {
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

		model, err := with(repository, func(model *todo.List) error {
			return model.Destroy(id)
		})
		if err != nil {
			handleError(w, err)
			return
		}

		vm := viewModel(model, r)
		render(w, r, templ, vm)
	})
}

func with(repository Repository, f func(list *todo.List) error) (*todo.List, error) {
	model, err := repository.FindList()
	if err != nil {
		return nil, err
	}

	if err := f(model); err != nil {
		return nil, err
	}

	if err := repository.SaveList(model); err != nil {
		return nil, err
	}
	return model, nil
}

func internalServerError(w http.ResponseWriter, err error) {
	log.Printf("Finder: %s", err.Error())
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

func badRequest(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusBadRequest)
}
