package segment

import (
	"context"
	"net/http"
	"time"

	"github.com/dezzerlol/avitotech-test-2023/pkg/payload"
)

// GetUserHistory godoc
// @Summary      Получение истории сегментов пользователя
// @Description  Метод получения истории сегментов пользователя за указанный месяц и год. На вход: год и месяц. На выходе ссылка на CSV файл.
// @Tags         Segment
// @Produce      json
// @Param        userId path string true "id пользователя"
// @Param        month query int true "месяц"
// @Param        year query int true "год"
// @Success      200  {object} object{report_link=string}
// @Failure      400,500  {object} object{error=string}
// @Router       /segment/history/{userId} [get]
func (h *handler) GetUserHistory(w http.ResponseWriter, r *http.Request) {
	userId, err := payload.ParamInt(r, "userId")
	if err != nil {
		payload.WriteJSON(w, http.StatusBadRequest, payload.Data{"error": err.Error()}, nil)
		return
	}

	month, err := payload.QueryInt(r, "month")
	if err != nil {
		payload.WriteJSON(w, http.StatusBadRequest, payload.Data{"error": err.Error()}, nil)
		return
	}

	year, err := payload.QueryInt(r, "year")
	if err != nil {
		payload.WriteJSON(w, http.StatusBadRequest, payload.Data{"error": err.Error()}, nil)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	downloadLink, err := h.segmentSvc.GetUserHistory(ctx, userId, month, year)

	if err != nil {
		payload.WriteJSON(w, http.StatusInternalServerError, payload.Data{"error": "Internal server error"}, nil)
		return
	}

	payload.WriteJSON(w, http.StatusOK, payload.Data{"report_link": downloadLink}, nil)
}
