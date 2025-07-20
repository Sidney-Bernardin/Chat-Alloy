package web

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	server "github.com/Sidney-Bernardin/Chat-Alloy/internal"
	"github.com/Sidney-Bernardin/Chat-Alloy/internal/service"
	"github.com/a-h/templ"
	"github.com/pkg/errors"
)

type Server struct {
	*http.Server

	Config *server.Config
	Logger *slog.Logger

	Service *service.Service
}

func (h *Server) Err(w io.Writer, r *http.Request, err error) {

}

func (h *Server) Respond(w io.Writer, r *http.Request, statusCode int, data any) {
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
		h.Err(w, r, errors.Wrap(err, "cannot write response"))
	}
}
