package handlers

import (
	"log/slog"
	"net/http"

	"github.com/ndrscodes/yasp/internal/util/templates"
)

type HomeHandler struct {
	registry     *templates.TemplateRegistry
	templateName string
}

func NewHomeHandler(reg *templates.TemplateRegistry, tn string) HomeHandler {
	return HomeHandler{
		registry:     reg,
		templateName: tn,
	}
}

func (h HomeHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	t, err := h.registry.Get()
	if err != nil {
		slog.Error("unable to obtain template", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, h.templateName, nil)
	if err != nil {
		slog.Error("template execution failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
