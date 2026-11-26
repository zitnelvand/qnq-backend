package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zitnelvand/qnq-backend/internal/config"
)

func main() {
	cfg, err := config.LoadHTTPConfig()
	if err != nil {
		log.Fatalf("qnq: load config: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	if err := http.ListenAndServe(cfg.Address(), r); err != nil {
		log.Fatalf("qnq: %v", err)
	}
}
