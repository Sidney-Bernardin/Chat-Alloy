package handlers

import (
	"net/http"
	"time"
)

type wrappedWritter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWritter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func (h *handler) mwLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &wrappedWritter{w, http.StatusOK}

		next.ServeHTTP(ww, r)
		h.log.Info("New request",
			"status", ww.statusCode,
			"method", r.Method,
			"path", r.URL.String(),
			"time", time.Since(start))
	})
}
