package httphelper

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Success")
	}))
	defer ts.Close()
	res, _ := Get(ts.URL, ts.Client())

	if res != "1313058fad6daf581fdcfea2e76ac3b1" {
		t.Fail()
	}
}
