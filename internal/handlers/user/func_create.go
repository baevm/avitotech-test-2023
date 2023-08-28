package user

import (
	"net/http"

	"github.com/dezzerlol/avitotech-test-2023/pkg/payload"
)

// Create godoc
// @Summary      Создание пользователя
// @Description  Метод создания пользователя.
// @Description  Используется в случае необходимости вручную добавить пользователя, так как при добавлении сегмента пользователь сохраняется автоматически.
// @Tags         User
// @Produce      json
// @Success      201  {object} object{user_id=int64}
// @Failure      400  {object} object{error=string}
// @Router       /user [post]
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	userId, err := h.userSvc.Create(r.Context())

	if err != nil {
		payload.WriteJSON(w, http.StatusInternalServerError, payload.Data{"error": "Internal server error"}, nil)
		return
	}

	payload.WriteJSON(w, http.StatusCreated, payload.Data{"user_id": userId}, nil)
}
