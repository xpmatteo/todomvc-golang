package web

import (
	"fmt"
	"github.com/xpmatteo/todomvc-golang/todo"
	"net/http"
)

type ViewModel map[string]interface{}

func viewModel(model *todo.List, r *http.Request) ViewModel {
	items := model.AllItems()
	path := determinePath(r)
	if path == pathCompleted {
		items = model.CompletedItems()
	} else if path == pathActive {
		items = model.ActiveItems()
	}
	return ViewModel{
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
