package web

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

type wrappedWritter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWritter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func (h *Server) MWLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &wrappedWritter{w, http.StatusOK}

		next.ServeHTTP(ww, r)

		h.Logger.Info("New request",
			"status", ww.statusCode,
			"method", r.Method,
			"path", r.URL.String(),
			"time", time.Since(start))
	})
}

func (h *Server) MWLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		http.SetCookie(w, &http.Cookie{
			Name:     "SESSION_ID",
			Value:    r.Context().Value("session_id").(uuid.UUID).String(),
			Domain:   h.Config.SESSION_COOKIE_DOMAIN,
			Expires:  time.Now().Add(h.Config.SESSION_DURATION),
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		})

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	})
}
