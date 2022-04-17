package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/bdharris08/scorekeeper"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const scoreTypeKey = "scoreType"

func (s *Server) HandleAddAction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		st, ok := r.Context().Value(scoreTypeKey).(string)
		if !ok {
			render.Render(w, r, ErrInternalServer(errors.New("failed to get scoreType from context")))
			return
		}

		fmt.Println(r.ContentLength)
		fmt.Println(r.Body)

		// validate the incoming json body
		body := make(map[string]interface{}, 2)
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			render.Render(w, r, ErrInternalServer(fmt.Errorf("reading request body: %w", err)))
			return
		}
		// we could validate specific fields here if we wanted
		marshaled, err := json.Marshal(body)
		if err != nil {
			render.Render(w, r, ErrInternalServer(fmt.Errorf("marshalling request body: %w", err)))
			return
		}

		payload := string(marshaled)

		if err := s.sk.AddAction(st, payload); err != nil {
			render.Render(w, r, ErrInternalServer(err))
			return
		}

		render.NoContent(w, r)
	}
}

// ScoreCtx gets the scoreType from the url parameter
func (s *Server) ScoreCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		scoreType := chi.URLParam(r, scoreTypeKey)
		// We need to check the scoreType for sql injection attacks
		// 	since `database/sql` doesn't let us use placeholders for table names.
		// A better architecture in scorekeeper wouldn't take table names as a parameter
		//	but that's some homework for another time.
		// Scorekeeper technically validates as well, but this lets us quit early.
		if ok := scorekeeper.ValidScoreType(s.sk, scoreType); !ok {
			render.Render(w, r, ErrInvalidRequest(fmt.Errorf("bad urlparam /scores/{scoreType}")))
			return
		}

		ctx := context.WithValue(r.Context(), scoreTypeKey, scoreType)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
