package segment

import (
	"net/http"

	"github.com/dezzerlol/avitotech-test-2023/internal/db/models"
	"github.com/dezzerlol/avitotech-test-2023/pkg/payload"
)

type DeleteRequest struct {
	Slug string `json:"slug"`
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	var req DeleteRequest

	if err := payload.ReadJSON(w, r, &req); err != nil {
		payload.WriteJSON(w, http.StatusBadRequest, payload.Data{"error": err.Error()}, nil)
		return
	}

	segment := &models.Segment{
		Slug: req.Slug,
	}

	if err := h.segment.DeleteBySlug(r.Context(), segment); err != nil {
		payload.WriteJSON(w, http.StatusBadRequest, payload.Data{"error": err.Error()}, nil)
		return
	}

	payload.WriteJSON(w, http.StatusOK, payload.Data{"message": "ok"}, nil)
}
