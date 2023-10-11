package server

import (
	"log/slog"
	"net/http"
	"reedsal/auth"
	"reedsal/messages"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	Router chi.Router
	DB     *sqlx.DB
	Redis  *redis.Client
}

func NewServer() *Server {
	return &Server{Router: chi.NewRouter()}
}

func (s *Server) WithMiddleware() *Server {
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(NewLoggingMiddleware())
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.Timeout(60 * time.Second))
	return s
}

func (s *Server) WithDB(db *sqlx.DB) *Server {
	s.DB = db
	return s
}

func (s *Server) WithRedis(client *redis.Client) *Server {
	s.Redis = client
	return s
}

func (s Server) Mount() Server {
	s.Router.Route("/reed/api", func(r chi.Router) { // TODO: Edit prefix
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"detail": "ok"}`))
		})

		r.Mount(
			"/auth",
			auth.NewAuthHandler(s.DB).Router(),
		)

		r.Mount(
			"/messages",
			messages.NewMessageHandler(s.DB, s.Redis).Router(),
		)
	})

	return s
}

func (s Server) Serve() {
	slog.Info("Starting to listen", "port", 8080)
	slog.Error(http.ListenAndServe(":8080", s.Router).Error())
}
