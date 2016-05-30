package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

var update = flag.String("golden", "", "update golden file")

func init() {
	flag.Parse()
}

func TestNewIndex(t *testing.T) {
	idx, err := NewIndex("testdata")
	if err != nil {
		t.Error(err)
	}

	for _, m := range idx {
		if m.Title == "" {
			t.Error("metadata title is empty")
		}
	}
}

func TestIndexHandler(t *testing.T) {
	idx, _ := NewIndex("testdata")
	ts := httptest.NewServer(http.HandlerFunc(idx.IndexHandler))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	golden := filepath.Join("testdata", "index.golden")
	if *update == "index" {
		ioutil.WriteFile(golden, body, 0644)
	}
	want, _ := ioutil.ReadFile(golden)
	if !bytes.Equal(body, want) {
		t.Errorf("\ngot:\n---\n%s---\nwant:\n---\n%s---", body, want)
	}
}
