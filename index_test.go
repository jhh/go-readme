package main

import (
	"bytes"
	"flag"
	"fmt"
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

	for c := range idx {
		for _, m := range idx[c] {
			if m.Title == "" {
				t.Error("metadata title is empty")
			}
		}
	}
}

func TestIndexHandler(t *testing.T) {
	idx, err := NewIndex("testdata")
	if err != nil {
		t.Error(err)
	}
	ts := httptest.NewServer(http.HandlerFunc(idx.IndexHandler))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	defer func() {
		if e := res.Body.Close(); e != nil {
			fmt.Print(e)
		}
	}()
	if err != nil {
		t.Fatal(err)
	}

	golden := filepath.Join("testdata", "index.golden")
	if *update == "index" {
		if err := ioutil.WriteFile(golden, body, 0644); err != nil {
			t.Fatal(err)
		}
	}
	want, _ := ioutil.ReadFile(golden)
	if !bytes.Equal(body, want) {
		t.Errorf("\ngot:\n---\n%s---\nwant:\n---\n%s---", body, want)
	}
}
