package segment

import (
	"net/http"

	"github.com/dezzerlol/avitotech-test-2023/pkg/payload"
)

type GetSegmentsForUserRequest struct {
	UserId int64 `json:"user_id"`
}

// GetSegmentsForUser godoc
// @Summary      Получение сегментов пользователя
// @Description Метод получения активных сегментов пользователя. Принимает на вход id пользователя.
// @Tags         segment
// @Accept       json
// @Produce      json
// @Param        body  body  GetSegmentsForUserRequest  true  "Данные пользователя"
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
