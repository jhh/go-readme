package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"text/template"

	"jhhgo.us/tmp/markdown"
)

var (
	port        = os.Getenv("PORT")
	contentTmpl = template.Must(
		template.ParseFiles(
			filepath.Join("templates", "content.html"),
			filepath.Join("templates", "layout.html"),
		),
	)
	validPath = regexp.MustCompile("^/(content)/([a-zA-Z0-9_]+\\.md)$")
)

func handleContent(w http.ResponseWriter, r *http.Request) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return
	}

	doc, err := markdown.NewDocument(filepath.Join(m[1], m[2]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	err = contentTmpl.ExecuteTemplate(w, "layout", doc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

func main() {
	log.SetFlags(0)

	http.HandleFunc("/content/", handleContent)

	idx, err := NewIndex("content")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", idx.IndexHandler)

	http.Handle("/css/", http.FileServer(http.Dir("static")))

	log.Println("listening on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
