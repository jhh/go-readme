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

// IndexEntry is an reference to a document.
type IndexEntry struct {
	Path string
	markdown.Meta
}

// Index holds the metadata for all content.
type Index []IndexEntry

// NewIndex loads the index.
func NewIndex(path string) (Index, error) {
	files, err := filepath.Glob(filepath.Join(path, "*.md"))
	if err != nil {
		return nil, err
	}

	idx := Index{}
	for _, f := range files {
		doc, err := markdown.NewDocument(f)
		if err != nil {
			return nil, err
		}
		idx = append(idx, IndexEntry{Path: filepath.ToSlash(f), Meta: doc.Meta})
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
