package segment

import (
	"net/http"

	"github.com/dezzerlol/avitotech-test-2023/pkg/payload"
)

// GetSegmentsForUser godoc
// @Summary      Получение сегментов пользователя
// @Description Метод получения активных сегментов пользователя. Принимает на вход id пользователя.
// @Tags         Segment
// @Produce      json
// @Param        userId path string true "id пользователя"
// @Success      200  {object} []models.Segment
// @Failure      400  {object} object{error=string}
// @Router       /segment/user/{userId} [get]
func (h *handler) GetSegmentsForUser(w http.ResponseWriter, r *http.Request) {
	userId, err := payload.ParamInt(r, "userId")
	if err != nil {
		payload.WriteJSON(w, http.StatusBadRequest, payload.Data{"error": err.Error()}, nil)
		return
	}

	segments, err := h.segmentSvc.GetUserSegments(r.Context(), userId)

	if err != nil {
		payload.WriteJSON(w, http.StatusInternalServerError, payload.Data{"error": "Internal server error"}, nil)
		return
	}

	payload.WriteJSON(w, http.StatusOK, payload.Data{"segments": segments}, nil)
}
