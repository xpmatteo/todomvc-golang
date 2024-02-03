package web

import (
	"html/template"
	"net/http"
	"xpug.it/todo"
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

func MakeDestroyHandler(list *todo.List) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			return
		}
		list.Destroy(todo.ItemId(r.PostForm.Get("todoItemId")))
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}
