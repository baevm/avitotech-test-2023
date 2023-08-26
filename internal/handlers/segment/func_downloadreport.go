package segment

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dezzerlol/avitotech-test-2023/pkg/payload"
	"github.com/go-chi/chi/v5"
)

// DownloadReport godoc
// @Summary      Скачивание отчета
// @Description  Метод скачивания csv отчета по истории сегментов пользователя.
// @Description  Отчет в формате: идентификатор пользователя 1;сегмент1;операция (добавление = 'I' / удаление = "D");дата и время
// @Tags         segment
// @Produce      text/csv
// @Param        fileName path string true "file_name.csv"
// @Success      200  {file} file
// @Failure      400  {object} object{error=string}
// @Header	 	 200 {string} Content-Type "text/csv"
// @Header 	 	 200 {string} Content-Disposition "attachment;filename=file_name"
// @Router       /segment/reports/{fileName} [get]
func (h *handler) DownloadReport(w http.ResponseWriter, r *http.Request) {
	fileName := chi.URLParam(r, "fileName")

	fullPath := fmt.Sprintf("./reports/%s", fileName)

	// Проверяем существует ли файл
	_, err := os.Stat(fullPath)

	if os.IsNotExist(err) {
		payload.WriteJSON(w, http.StatusBadRequest, payload.Data{"error": "file not found"}, nil)
		return
	}

	// Устанавливаем заголовки позволяющие браузеру скачать файл
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename="+fileName)

	// TODO: добавить удаление файла после скачивания
	http.ServeFile(w, r, fullPath)
}
