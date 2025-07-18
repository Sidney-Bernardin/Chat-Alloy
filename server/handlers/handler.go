package handlers

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/Sidney-Bernardin/Chat-Alloy/server"
	"github.com/Sidney-Bernardin/Chat-Alloy/server/service"
	"github.com/a-h/templ"
	"github.com/pkg/errors"
)

type handler struct {
	*http.Server

	cfg *server.Config
	log *slog.Logger

	svc *service.Service
}

func New(cfg *server.Config, log *slog.Logger, svc *service.Service) *handler {
	h := &handler{
		&http.Server{
			Addr: cfg.ADDR,
		},
		cfg, log, svc,
	}

	r := http.NewServeMux()
	r.Handle("/", newPagesHandler(h))
	h.Server.Handler = h.mwLog(r)

	return h
}

func (h *handler) err(w io.Writer, r *http.Request, err error) {

}

func (h *handler) respond(w io.Writer, r *http.Request, statusCode int, data any) {
	if w, ok := w.(http.ResponseWriter); ok {
		w.WriteHeader(statusCode)
	}

	var err error
	switch d := data.(type) {
	case templ.Component:
		err = d.Render(r.Context(), w)
	default:
		err = json.NewEncoder(w).Encode(data)
	}

	if err != nil {
		h.err(w, r, errors.Wrap(err, "cannot write response"))
	}
}
