package server

import (
	"fmt"
	"net/http"

	"github.com/bdharris08/scorekeeper"
	"github.com/go-chi/chi/v5"
)

// handle http requests and response
// call scorekeeper
// two request types
// - add action
// - get stats

// Following along with https://github.com/go-chi/chi
// based on advice from: https://www.veritone.com/blog/how-i-write-go-http-services-after-seven-years/

type Server struct {
	sk     *scorekeeper.ScoreKeeper
	router chi.Router
}

func NewServer(sk *scorekeeper.ScoreKeeper) *Server {
	s := &Server{
		sk: sk,
	}

	s.router = chi.NewRouter()
	s.routes()

	return s
}

func (s *Server) ListenAndServe(address string) {
	fmt.Println("Listening on ", address)
	http.ListenAndServe(address, s.router)
}

// TODO
func (s *Server) HandleAddAction() http.HandlerFunc {
	return nil
}

// TODO
func (s *Server) HandleGetStats() http.HandlerFunc {
	return nil
}
