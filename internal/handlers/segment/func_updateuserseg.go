package segment

import (
	"net/http"

	"github.com/dezzerlol/avitotech-test-2023/pkg/payload"
)

type UpdateUserSegmentsRequest struct {
	UserId         int64    `json:"user_id" validate:"required,min=1" example:"1"`
	AddSegments    []string `json:"add_segments" validate:"dive,min=3" example:"AVITO_VOICE_MESSAGES,AVITO_DISCOUNT_50"`
	TTL            int64    `json:"ttl" validate:"omitempty,min=1" example:"1000"`
	DeleteSegments []string `json:"delete_segments" validate:"dive,min=3" example:"AVITO_DISCOUNT_10"`
}

// UpdateUserSegments godoc
// @Summary      Добавление/удаление сегментов у пользователя
// @Description  Метод добавления пользователя в сегмент. Принимает массив slug (названий) сегментов которые нужно добавить пользователю,
// @Description  массив slug (названий) сегментов которые нужно удалить у пользователя, id пользователя, ttl (в секундах).
// @Tags         Segment
// @Accept       json
// @Produce      json
// @Param        body  body  UpdateUserSegmentsRequest  true  "Данные сегмента и пользователя"
// @Success      200  {object} object{segments_added=int,segments_deleted=int}
// @Failure      400  {object} object{error=string}
// @Router       /segment/user [post]
func (h *handler) UpdateUserSegments(w http.ResponseWriter, r *http.Request) {
	var req UpdateUserSegmentsRequest

	if err := payload.ReadJSON(w, r, &req); err != nil {
		payload.WriteJSON(w, http.StatusBadRequest, payload.Data{"error": err.Error()}, nil)
		return
	}

	if errs := payload.Validate(req); errs != nil {
		payload.WriteJSON(w, http.StatusBadRequest, payload.Data{"error": errs}, nil)
		return
	}

	segmentsAdded, segmentsDeleted, err := h.segmentSvc.UpdateUserSegments(
		r.Context(),
		req.UserId,
		req.AddSegments,
		req.TTL,
		req.DeleteSegments,
	)

	if err != nil {
		payload.WriteJSON(w, http.StatusBadRequest, payload.Data{"error": err.Error()}, nil)
		return
	}

	payload.WriteJSON(w, http.StatusOK, payload.Data{"segments_added": segmentsAdded, "segments_deleted": segmentsDeleted}, nil)
}
