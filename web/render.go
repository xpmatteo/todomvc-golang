package web

import (
	"encoding/json"
	"html/template"
	"net/http"
)

const (
	headerKeyAccept = "accept"
)

func render(w http.ResponseWriter, r *http.Request, templ *template.Template, model ViewModel) {
	if acceptsJson(r) {
		w.Header().Add("content-type", "application/json")
		err := json.NewEncoder(w).Encode(model)
		if err != nil {
			panic(err.Error())
		}
		return
	}
	err := templ.Execute(w, model)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func acceptsJson(r *http.Request) bool {
	accept := r.Header.Get(headerKeyAccept)

	return startsWith("application/json", accept)
}

func startsWith(prefix string, s string) bool {
	lenPrefix := len(prefix)

	return lenPrefix <= len(s) && s[:lenPrefix] == prefix
}
