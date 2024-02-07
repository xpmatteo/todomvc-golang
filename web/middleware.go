package web

import (
	"log"
	"net/http"
	"time"
)

func GET(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

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

func Slowdown(delayMilli int, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(delayMilli) * time.Millisecond)
		handler.ServeHTTP(w, r)
	})
}
