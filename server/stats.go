package server

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

// HandleGetStats returns average scores from the ScoreKeeper
func (s *Server) HandleGetStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		st, ok := r.Context().Value(scoreTypeKey).(string)
		if !ok {
			render.Render(w, r, ErrInternalServer(errors.New("failed to get scoreType from context")))
			return
		}

		stats, err := s.sk.GetStats(st)
		if err != nil {
			render.Render(w, r, ErrInternalServer(err))
			return
		}

		// Stats is already json-encoded string so no need to do any rendering
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(stats))
	}
}
