package main

import (
	"html/template"
	"net/http"
	"path/filepath"

	"jhhgo.us/markdown"
)

var indexTmpl = template.Must(
	template.ParseFiles(
		filepath.Join("templates", "index.html"),
		filepath.Join("templates", "layout.html"),
	),
)

// Index holds the metadata for all content.
type Index []markdown.Meta

// NewIndex loads the index.
func NewIndex(path string) (Index, error) {
	files, err := filepath.Glob(filepath.Join(path, "*.md"))
	if err != nil {
		return nil, err
	}

	return makeIndex(files)
}

func makeIndex(files []string) (Index, error) {
	idx := make(Index, 0)
	for _, f := range files {
		doc, err := markdown.NewDocument(f)
		if err != nil {
			return nil, err
		}
		idx = append(idx, doc.Meta)
	}
	return idx, nil
}

// IndexHandler displays the page index.
func (idx Index) IndexHandler(w http.ResponseWriter, r *http.Request) {
	err := indexTmpl.ExecuteTemplate(w, "layout", idx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
