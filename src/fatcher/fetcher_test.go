package fatcher

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetch(t *testing.T) {
	const content = "this is test content!"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(content))
	}))
	defer ts.Close()
	body, err := Fetch(ts.URL)
	if err != nil {
		t.Error(err)
	}
	if string(body) != content {
		t.Errorf("Content except %s,actual %s",
			content, string(body))
	}
}
