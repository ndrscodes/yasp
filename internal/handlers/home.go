package handlers

import (
	"errors"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/ndrscodes/yasp/internal/templates"
)

type HomeHandler struct {
	template *template.Template
}

func NewHomeHandler(root *template.Template) (HomeHandler, error) {
	if root == nil {
		return HomeHandler{}, errors.New("root is nil")
	}

	t, err := root.ParseFS(templates.Files, "pages/home/*")
	if err != nil {
		return HomeHandler{}, err
	}

	return HomeHandler{
		template: t,
	}, err
}

func (h *HomeHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	err := h.template.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		slog.Error("template execution failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
