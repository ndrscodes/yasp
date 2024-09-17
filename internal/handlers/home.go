package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/ndrscodes/yasp/internal/db"
	templates "github.com/ndrscodes/yasp/internal/web/templates/pages/home"
)

type HomeHandler struct {
	systemRepository *db.SystemRepository
}

func NewHomeHandler(repo *db.SystemRepository) (HomeHandler, error) {
	if repo == nil {
		return HomeHandler{}, errors.New("systemRepository is nil")
	}

	return HomeHandler{
		systemRepository: repo,
	}, nil
}

func (h *HomeHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	systems, err := h.systemRepository.GetAll()

	err = templates.Index(templates.SystemData{Systems: systems, Error: err}).Render(r.Context(), w)

	if err != nil {
		slog.Error("template execution failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
