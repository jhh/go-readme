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
type Index map[string][]markdown.Meta

// NewIndex loads the index.
func NewIndex(path string) (Index, error) {
	// will read all docs in content directory given in path
	files, err := filepath.Glob(filepath.Join(path, "*.md"))
	if err != nil {
		return nil, err
	}

	// get metadata for each document
	meta, err := collectMeta(files)
	if err != nil {
		return nil, err
	}

	// group by categories
	idx := make(Index)
	for _, m := range meta {
		for _, c := range m.Categories {
			idx[c] = append(idx[c], m)
		}
	}

	return idx, err
}

func collectMeta(files []string) ([]markdown.Meta, error) {
	var idx []markdown.Meta
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
