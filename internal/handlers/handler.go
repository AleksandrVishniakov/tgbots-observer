package handlers

import (
	"log/slog"
	"net/http"

	"github.com/AleksandrVishniakov/tgbots-util/http/middlewares"
)

type Handler struct {
	log *slog.Logger
}

func New(log *slog.Logger) *Handler {
	return &Handler{
		log: log,
	}
}

func (h *Handler) InitRoutes(routers map[string]http.Handler) http.Handler {
	recovery := middlewares.Recovery(h.log)
	logger := middlewares.Logger(h.log)

	mux := http.NewServeMux()

	for prefix, router := range routers {
		mux.Handle(prefix + "/", http.StripPrefix(prefix, router))
	}

	mux.HandleFunc("/ping", h.Ping)

	return recovery(logger(mux))
}

func (h *Handler) Ping(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
