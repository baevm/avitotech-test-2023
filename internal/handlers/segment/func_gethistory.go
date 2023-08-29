package segment

import (
	"context"
	"fmt"
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
// @Success      200  {object} object{report=string}
// @Failure      400  {object} object{error=string}
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

	filePath, err := h.segmentSvc.GetUserHistory(ctx, userId, month, year)

	if err != nil {
		payload.WriteJSON(w, http.StatusInternalServerError, payload.Data{"error": "Internal server error"}, nil)
		return
	}

	// Ссылка для скачивания файла в формате addr:port/segment/reports/file_name.csv
	addr := fmt.Sprintf("%s:%s", h.config.REPORTS_HOST, h.config.API_PORT)
	downloadLink := addr + "/segment" + filePath

	payload.WriteJSON(w, http.StatusOK, payload.Data{"report": downloadLink}, nil)
}
