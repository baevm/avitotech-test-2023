package segment

import (
	"net/http"

	"github.com/dezzerlol/avitotech-test-2023/internal/db/models"
	"github.com/dezzerlol/avitotech-test-2023/pkg/payload"
)

type DeleteRequest struct {
	Slug string `json:"slug" validate:"required,min=3" example:"AVITO_VOICE_MESSAGES"`
}

// Delete godoc
// @Summary      Удаление сегмента
// @Description  Метод удаления сегмента. Принимает slug (название) сегмента.
// @Tags         Segment
// @Accept       json
// @Produce      json
// @Param        body  body  DeleteRequest  true  "Данные сегмента"
// @Success      200  {object} object{message=string}
// @Failure      400  {object} object{error=string}
// @Router       /segment [delete]
func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	var req DeleteRequest

	if err := payload.ReadJSON(w, r, &req); err != nil {
		payload.WriteJSON(w, http.StatusBadRequest, payload.Data{"error": err.Error()}, nil)
		return
	}

	if errs := payload.Validate(req); errs != nil {
		payload.WriteJSON(w, http.StatusBadRequest, payload.Data{"error": errs}, nil)
		return
	}

	segment := &models.Segment{
		Slug: req.Slug,
	}

	if err := h.segmentSvc.DeleteBySlug(r.Context(), segment); err != nil {
		payload.WriteJSON(w, http.StatusBadRequest, payload.Data{"error": err.Error()}, nil)
		return
	}

	payload.WriteJSON(w, http.StatusOK, payload.Data{"message": "ok"}, nil)
}
