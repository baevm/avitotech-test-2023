package segment

import (
	"net/http"

	"github.com/dezzerlol/avitotech-test-2023/pkg/payload"
)

type GetSegmentsForUserRequest struct {
	UserId int64 `json:"user_id"`
}

// GetSegmentsForUser godoc
// @Summary      Get user segments
// @Description  Получение всех сегментов для пользователя
// @Tags         segment
// @Accept       json
// @Produce      json
// @Param        body  body  GetSegmentsForUserRequest  true  "User data"
// @Success      200  {object} []models.Segment
// @Failure      400  {object} object{error=string}
// @Router       /segment/user [get]
func (h *handler) GetSegmentsForUser(w http.ResponseWriter, r *http.Request) {
	var req GetSegmentsForUserRequest

	if err := payload.ReadJSON(w, r, &req); err != nil {
		payload.WriteJSON(w, http.StatusBadRequest, payload.Data{"error": err.Error()}, nil)
		return
	}

	segments, err := h.segment.GetUserSegments(r.Context(), req.UserId)

	if err != nil {
		payload.WriteJSON(w, http.StatusBadRequest, payload.Data{"error": err.Error()}, nil)
		return
	}

	payload.WriteJSON(w, http.StatusOK, payload.Data{"segments": segments}, nil)
}
