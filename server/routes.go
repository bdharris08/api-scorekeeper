package server

import (
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) routes() {
	// middleware stack (recommended in chi docs)
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)

	s.router.Use(middleware.Timeout(60 * time.Second))

	s.router.Get("/", s.handleHello())
}
