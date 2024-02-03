package web

import (
	"html/template"
	"net/http"
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
		err := templ.Execute(w, model)
		if err != nil {
			panic(err.Error())
		}
	})
}
