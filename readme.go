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
	port = os.Getenv("PORT")
)

func main() {
	for p := range urlMap {
		http.HandleFunc(p, func(w http.ResponseWriter, r *http.Request) {
			fn := filepath.Join("content", urlMap[p])
			doc, err := markdown.NewDocument(fn)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			err = templates.ExecuteTemplate(w, "view.html", doc)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
	}
	fmt.Println("Listening on http://localhost:"+port)
	http.ListenAndServe(":"+port, nil)
}

var urlMap = map[string]string{
	"/rpi-ruby": "20160527_rpi_ruby.md",
}
