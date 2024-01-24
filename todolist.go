//go:build ignore

package main

import (
	// "fmt"
	"html/template"
	"log"
	"net/http"
)

const port = "8080"

var templates = template.Must(template.ParseFiles("index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
	// fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", indexHandler)

	fs := http.FileServer(http.Dir("./public/css"))
	http.Handle("/css", fs)

	log.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
