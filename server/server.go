package server

import (
	"fmt"
	"net/http"

	"github.com/bdharris08/scorekeeper"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/docgen"
)

// Server facilitates http requests to a ScoreKeeper
// It uses the chi router
// The pattern is adapted from:
// - https://github.com/go-chi/chi
// - https://www.veritone.com/blog/how-i-write-go-http-services-after-seven-years/
type Server struct {
	sk     *scorekeeper.ScoreKeeper
	router chi.Router
}

// NewServer returns a new initialized server
// the provided scorekeeper should:
// - be running
// - have an initialized Store
func NewServer(sk *scorekeeper.ScoreKeeper) *Server {
	s := &Server{
		sk: sk,
	}

	s.router = chi.NewRouter()
	s.routes()

	return s
}

// Usage returns the Server routes
func (s *Server) Usage() string {
	return docgen.JSONRoutesDoc(s.router)
}

// ListenAndServe starts the Server listening for requests
func (s *Server) ListenAndServe(address string) {
	fmt.Println("Listening on ", address)
	http.ListenAndServe(address, s.router)
}
