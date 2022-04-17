package server

import (
	"net/http"

	"github.com/go-chi/render"
)

func (s *Server) handleHello() http.HandlerFunc {
	// keep handler request and response types in each handler,
	// only need to export shared response types
	type response struct {
		Greeting string `json:"greeting"`
	}

	// everything initialized in this scope is only done once
	// the returned handlerFunc has access to these through closures
	standard := response{
		Greeting: "Hello World!",
	}

	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, standard)
	}
}
