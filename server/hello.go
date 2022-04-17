package server

import (
	"encoding/json"
	"net/http"
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
		json, err := json.Marshal(standard)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		_, err = w.Write(json)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
	}
}
