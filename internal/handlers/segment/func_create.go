package segment

import (
	"net/http"

	"github.com/dezzerlol/avitotech-test-2023/internal/db/models"
	"github.com/dezzerlol/avitotech-test-2023/pkg/payload"
)

type CreateRequest struct {
	Slug string `json:"slug"`
}

// Create godoc
// @Summary      Создание сегмента
// @Description  Метод создания сегмента. Принимает slug (название) сегмента.
// @Tags         segment
// @Accept       json
// @Produce      json
// @Param        body  body  CreateRequest  true  "Данные сегмента"
// @Success      201  {object} object{created_at=string}
// @Failure      400  {object} object{error=string}
// @Router       /segment [post]
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

	payload.WriteJSON(w, http.StatusCreated, payload.Data{"created_at": segment.CreatedAt}, nil)
}
