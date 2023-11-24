package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mock_user "github.com/dezzerlol/avitotech-test-2023/internal/handlers/user/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_CreateUser(t *testing.T) {
	t.Run("Should successfully create user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSegmentSvc := mock_user.NewMockUserService(ctrl)

		mockSegmentSvc.EXPECT().Create(gomock.Any()).Return(int64(1), nil)

		handler := NewHandler(nil, mockSegmentSvc)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/user", nil)
		handler.Create(w, r)

		require.Equal(t, http.StatusCreated, w.Code)
		require.JSONEq(t, `{"user_id":1}`, w.Body.String())
	})
}
