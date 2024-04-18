package handlers

import (
	"log/slog"
	"net/http"

	"github.com/ndrscodes/yasp/internal/util/templates"
)

type HomeHandler struct {
	registry *templates.TemplateRegistry
	page     string
}

func NewHomeHandler(reg *templates.TemplateRegistry, p string) HomeHandler {
	_, err := reg.Register(p)
	if err != nil {
		slog.Error("unable to set up template registry", "error", err)
		panic(err)
	}

	return HomeHandler{
		registry: reg,
		page:     p,
	}
}

func (h HomeHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	t, err := h.registry.Get(h.page)
	if err != nil {
		slog.Error("unable to obtain template", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		slog.Error("template execution failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
