package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bdharris08/scorekeeper"
	"github.com/bdharris08/scorekeeper/score"
	"github.com/bdharris08/scorekeeper/store"

	// using MatRyer's cool `is` testing package
	"github.com/matryer/is"
)

func TestGetStats(t *testing.T) {
	is := is.New(t)
	s := &store.MemoryStore{
		S: map[string]map[string][]score.Score{
			"test": {
				"test": []score.Score{
					&score.TestScore{
						TName:  "hop",
						TValue: float64(1),
					},
				},
			},
		},
	}
	f := score.ScoreFactory{
		"test": func() score.Score { return &score.TestScore{} },
	}
	sk, err := scorekeeper.New(s, f)
	if err != nil {
		t.Fatal(err)
	}

	sk.Start()
	defer sk.Stop()
	srv := NewServer(sk)

	req, err := http.NewRequest("GET", "/scores/test/average", nil)
	is.NoErr(err)

	w := httptest.NewRecorder()
	srv.router.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	is.Equal(resp.StatusCode, http.StatusOK)
	buf := make([]byte, w.Body.Len())
	_, err = resp.Body.Read(buf)
	is.NoErr(err)

	expected := `[{"action":"test","avg":1}]`
	is.Equal(string(buf), expected)
}

func TestGetStatsEmpty(t *testing.T) {
	is := is.New(t)
	sk := setupScoreKeeper(t)
	sk.Start()
	defer sk.Stop()
	srv := NewServer(sk)

	req, err := http.NewRequest("GET", "/scores/test/average", nil)
	is.NoErr(err)

	w := httptest.NewRecorder()
	srv.router.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	is.Equal(resp.StatusCode, http.StatusInternalServerError)
	buf := make([]byte, w.Body.Len())
	_, err = resp.Body.Read(buf)
	is.NoErr(err)

	expected := "no scores found"
	is.True(strings.Contains(string(buf), expected))
}

func TestGetUnknownStat(t *testing.T) {
	is := is.New(t)
	sk := setupScoreKeeper(t)
	sk.Start()
	defer sk.Stop()
	srv := NewServer(sk)

	req, err := http.NewRequest("GET", "/scores/test/median", nil)
	is.NoErr(err)

	w := httptest.NewRecorder()
	srv.router.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	is.Equal(resp.StatusCode, http.StatusNotFound)
	fmt.Println(resp.Body)
}
