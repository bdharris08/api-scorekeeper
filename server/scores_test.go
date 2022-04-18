package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bdharris08/scorekeeper"
	"github.com/bdharris08/scorekeeper/score"
	"github.com/bdharris08/scorekeeper/store"

	// using MatRyer's cool `is` testing package
	"github.com/matryer/is"
)

// setupScoreKeeper is a helper function to set up tests
func setupScoreKeeper(t *testing.T) *scorekeeper.ScoreKeeper {
	s := &store.MemoryStore{}
	f := score.ScoreFactory{
		"test": func() score.Score { return &score.TestScore{} },
	}
	sk, err := scorekeeper.New(s, f)
	if err != nil {
		t.Fatal(err)
	}
	return sk
}

func TestHandleAddActionBadType(t *testing.T) {
	is := is.New(t)
	sk := setupScoreKeeper(t)
	sk.Start()
	defer sk.Stop()
	srv := NewServer(sk)

	req, err := http.NewRequest("POST", "/scores/baad", nil)
	is.NoErr(err)

	w := httptest.NewRecorder()
	srv.router.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	is.Equal(resp.StatusCode, http.StatusBadRequest)
	is.Equal(resp.Status, "400 Bad Request")
}

func TestHandleAddAction(t *testing.T) {
	is := is.New(t)
	sk := setupScoreKeeper(t)
	sk.Start()
	defer sk.Stop()
	srv := NewServer(sk)

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(score.TestScore{
		TName:  "t",
		TValue: float64(1),
	})
	is.NoErr(err)

	req, err := http.NewRequest("POST", "/scores/test", &buf)
	is.NoErr(err)

	w := httptest.NewRecorder()
	srv.router.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	is.Equal(resp.StatusCode, http.StatusNoContent)
}
