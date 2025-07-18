package handlers

import (
	"net/http"

	"github.com/Sidney-Bernardin/Chat-Alloy/web/pages/home"
)

type PagesHandler struct {
	*handler
	router *http.ServeMux
}

func newPagesHandler(h *handler) *PagesHandler {
	pages := &PagesHandler{h, http.NewServeMux()}

	pages.router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("tmp/dist"))))
	pages.router.HandleFunc("/", pages.handleHome)

	return pages
}

func (h *PagesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *PagesHandler) handleHome(w http.ResponseWriter, r *http.Request) {
	h.respond(w, r, http.StatusOK, home.Home())
}
