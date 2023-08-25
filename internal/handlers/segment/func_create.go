package segment

import (
	"net/http"

	"github.com/dezzerlol/avitotech-test-2023/internal/db/models"
	"github.com/dezzerlol/avitotech-test-2023/pkg/payload"
)

type CreateRequest struct {
	Slug string `json:"slug"`
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest

	if err := payload.ReadJSON(w, r, &req); err != nil {
		payload.WriteJSON(w, http.StatusBadRequest, payload.Data{"error": err.Error()}, nil)
		return
	}

	segment := &models.Segment{
		Slug: req.Slug,
	}

	if err := h.segment.Create(r.Context(), segment); err != nil {
		payload.WriteJSON(w, http.StatusBadRequest, payload.Data{"error": err.Error()}, nil)
		return
	}

	payload.WriteJSON(w, http.StatusOK, payload.Data{"created_at": segment.CreatedAt}, nil)
}
