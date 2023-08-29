package segment

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/dezzerlol/avitotech-test-2023/internal/db/models"
	"github.com/dezzerlol/avitotech-test-2023/internal/repo"
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

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err := h.segmentSvc.DeleteBySlug(ctx, segment)

	if err != nil {
		switch {
		case errors.Is(err, repo.ErrSegmentNotFound):
			payload.WriteJSON(w, http.StatusNotFound, payload.Data{"error": err.Error()}, nil)
			return
		default:
			payload.WriteJSON(w, http.StatusInternalServerError, payload.Data{"error": "Internal server error"}, nil)
			return
		}
	}

	payload.WriteJSON(w, http.StatusOK, payload.Data{"message": "ok"}, nil)
}
