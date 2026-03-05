package app

import (
	"net/http"
	"time"

	"github.com/IbnBaqqi/gitchat/internal/config"
	"github.com/IbnBaqqi/gitchat/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(cfg *config.Config, db *database.DB) http.Handler {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	return r
}