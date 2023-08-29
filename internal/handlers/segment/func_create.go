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

type CreateRequest struct {
	Slug        string `json:"slug" validate:"required,min=3" example:"AVITO_VOICE_MESSAGES"`
	UserPercent int8   `json:"user_percent" validate:"omitempty,min=1,max=100" example:"50"`
}

// Create godoc
// @Summary      Создание сегмента
// @Description  Метод создания сегмента. Принимает slug (название) сегмента.
// @Description  Если указан user_percent, то сегмент будет добавлен проценту от общего числа случайным пользователям.
// @Tags         Segment
// @Accept       json
// @Produce      json
// @Param        body  body  CreateRequest  true  "Запрос на создание"
// @Success      201  {object} object{created_at=string}
// @Failure      400  {object} object{error=string}
// @Router       /segment [post]
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest

	if err := payload.ReadJSON(w, r, &req); err != nil {
		payload.WriteJSON(w, http.StatusBadRequest, payload.Data{"error": err.Error()}, nil)
		return
	}

	if errs := payload.Validate(req); errs != nil {
		payload.WriteJSON(w, http.StatusBadRequest, payload.Data{"error": errs}, nil)
		return
	}

	segment := &models.Segment{
		Slug:        req.Slug,
		UserPercent: req.UserPercent,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err := h.segmentSvc.Create(ctx, segment)

	if err != nil {
		switch {
		case errors.Is(err, repo.ErrSegmentAlreadyExists):
			payload.WriteJSON(w, http.StatusBadRequest, payload.Data{"error": err.Error()}, nil)
			return
		default:
			payload.WriteJSON(w, http.StatusInternalServerError, payload.Data{"error": "Internal server error"}, nil)
			return
		}
	}

	payload.WriteJSON(w, http.StatusCreated, payload.Data{"created_at": segment.CreatedAt}, nil)
}
