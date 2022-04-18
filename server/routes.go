package server

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// routes sets up the chi routes and middleware
func (s *Server) routes() {
	// middleware stack (recommended in chi docs)
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)

	s.router.Use(middleware.Timeout(60 * time.Second))

	s.router.Route("/hello", func(r chi.Router) {
		r.Get("/", s.HandleHello())
	})

	s.router.Route("/scores/{scoreType}", func(r chi.Router) {
		r.Use(s.ScoreCtx)
		r.Post("/", s.HandleAddAction())

		r.Route("/average", func(r chi.Router) {
			r.Get("/", s.HandleGetStats())
		})
	})
}
