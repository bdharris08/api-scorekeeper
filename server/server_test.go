package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	// using MatRyer's cool `is` testing package
	"github.com/matryer/is"
)

func TestHandleHello(t *testing.T) {
	is := is.New(t)
	srv := NewServer(nil)

	req, err := http.NewRequest("GET", "/hello", nil)
	is.NoErr(err)

	w := httptest.NewRecorder()
	srv.router.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	is.Equal(resp.StatusCode, http.StatusOK)
	buf := make([]byte, w.Body.Len())
	_, err = resp.Body.Read(buf)
	is.NoErr(err)

	expected := "{\"greeting\":\"Hello World!\"}\n"
	is.Equal(string(buf), expected)
}
