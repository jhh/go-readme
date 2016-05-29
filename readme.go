package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"jhhgo.us/tmp/markdown"
)

var (
	templates = template.Must(template.ParseGlob("templates/*.html"))
	port      = os.Getenv("PORT")
)

func handleContent(filename string, w http.ResponseWriter, r *http.Request) {
	doc, err := markdown.NewDocument(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = templates.ExecuteTemplate(w, "view.html", doc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	for p := range urlMap {
		http.HandleFunc(p, func(w http.ResponseWriter, r *http.Request) {
			fn := filepath.Join("content", urlMap[p])
			handleContent(fn, w, r)
		})
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fn := filepath.Join("content", "index.md")
		handleContent(fn, w, r)
	})
	fmt.Println("Listening on http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}

var urlMap = map[string]string{
	"/rpi-ruby":      "20160527_rpi_ruby.md",
	"/deploy-readme": "20160528_foreman.md",
}
