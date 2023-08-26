package segment

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dezzerlol/avitotech-test-2023/pkg/payload"
)

type HistoryRequest struct {
	UserId int64 `json:"user_id"`
	Month  int   `json:"month"`
	Year   int   `json:"year"`
}

// GetUserHistory godoc
// @Summary      Получение истории сегментов пользователя
// @Description  Метод получения истории сегментов пользователя за указанный месяц и год. На вход: год и месяц. На выходе ссылка на CSV файл.
// @Tags         segment
// @Accept       json
// @Produce      json
// @Param        body  body  HistoryRequest  true  "Данные сегмента"
// @Success      201  {object} object{report=string}
// @Failure      400  {object} object{error=string}
// @Router       /segment/user/history [get]
func (h *handler) GetUserHistory(w http.ResponseWriter, r *http.Request) {
	var req HistoryRequest

	if err := payload.ReadJSON(w, r, &req); err != nil {
		payload.WriteJSON(w, http.StatusBadRequest, payload.Data{"error": err.Error()}, nil)
		return
	}

	date := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)

	filePath, err := h.segment.GetUserHistory(r.Context(), req.UserId, date)

	if err != nil {
		payload.WriteJSON(w, http.StatusBadRequest, payload.Data{"error": err.Error()}, nil)
		return
	}

	// Ссылка для скачивания файла в формате addr:port/segments/reports/file_name
	addr := fmt.Sprintf("%s:%s", h.config.HTTP_HOST, h.config.HTTP_PORT)
	downloadLink := addr + "/segments" + filePath

	payload.WriteJSON(w, http.StatusCreated, payload.Data{"report": downloadLink}, nil)
}
