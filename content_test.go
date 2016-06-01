package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func TestContentHandler(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleContent(w, r)
	}))
	defer ts.Close()

	// override validPath set in readme.go
	saveValidPath := validPath
	validPath = regexp.MustCompile("^/(testdata)/(ok_cat_bar1\\.md)$")
	res, err := http.Get(ts.URL + "/testdata/ok_cat_bar1.md")
	if err != nil {
		t.Fatal(err)
	}
	// restore validPath
	validPath = saveValidPath

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	want := "<title>Category bar: 2016-04-27T23:45:46.000Z</title>"
	if !strings.Contains(string(body), want) {
		t.Errorf("body does not contain %q, got:\n%q\n", want, body)
	}
	want = "<h1>Category bar: 2016-04-27T23:45:46.000Z</h1>"
	if !strings.Contains(string(body), want) {
		t.Errorf("body does not contain %q, got:\n%q\n", want, body)
	}
}
