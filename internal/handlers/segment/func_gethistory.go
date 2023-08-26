package segment

import (
	"fmt"
	"net/http"

	"github.com/dezzerlol/avitotech-test-2023/pkg/payload"
)

// GetUserHistory godoc
// @Summary      Получение истории сегментов пользователя
// @Description  Метод получения истории сегментов пользователя за указанный месяц и год. На вход: год и месяц. На выходе ссылка на CSV файл.
// @Tags         segment
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

	filePath, err := h.segment.GetUserHistory(r.Context(), userId, month, year)

	if err != nil {
		payload.WriteJSON(w, http.StatusBadRequest, payload.Data{"error": err.Error()}, nil)
		return
	}

	// Ссылка для скачивания файла в формате addr:port/segments/reports/file_name
	addr := fmt.Sprintf("%s:%s", h.config.HTTP_HOST, h.config.HTTP_PORT)
	downloadLink := addr + "/segments" + filePath

	payload.WriteJSON(w, http.StatusOK, payload.Data{"report": downloadLink}, nil)
}
