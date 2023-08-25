package segment

import (
	"net/http"

	"github.com/dezzerlol/avitotech-test-2023/pkg/payload"
)

type UpdateUserSegmentsRequest struct {
	UserId         int64    `json:"user_id"`
	AddSegments    []string `json:"add_segments"`
	DeleteSegments []string `json:"delete_segments"`
}

func (h *handler) UpdateUserSegments(w http.ResponseWriter, r *http.Request) {
	var req UpdateUserSegmentsRequest

	if err := payload.ReadJSON(w, r, &req); err != nil {
		payload.WriteJSON(w, http.StatusBadRequest, payload.Data{"error": err.Error()}, nil)
		return
	}

	segmentsAdded, segmentsDeleted, err := h.segment.UpdateUserSegments(r.Context(), req.UserId, req.AddSegments, req.DeleteSegments)

	if err != nil {
		payload.WriteJSON(w, http.StatusBadRequest, payload.Data{"error": err.Error()}, nil)
		return
	}

	payload.WriteJSON(w, http.StatusOK, payload.Data{"segments_added": segmentsAdded, "segments_deleted": segmentsDeleted}, nil)
}
