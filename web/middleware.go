package web

import (
	"log"
	"net/http"
)

func Logging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)

		var body interface{}
		if len(r.Form) > 0 {
			body = r.Form
		} else {
			body = ""
		}
		log.Printf("%-4s %s %s", r.Method, r.URL.String(), body)
	})
}
