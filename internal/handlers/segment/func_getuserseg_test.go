package segment

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dezzerlol/avitotech-test-2023/internal/db/models"
	mock_segment "github.com/dezzerlol/avitotech-test-2023/internal/handlers/segment/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_GetUserSegments(t *testing.T) {
	t.Run("Should return 200 and user segments", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSegmentSvc := mock_segment.NewMockSegmentService(ctrl)
		handler := NewHandler(nil, mockSegmentSvc)

		var userId int64 = 1

		segments := make([]*models.Segment, 3)

		mockSegmentSvc.EXPECT().GetUserSegments(gomock.Any(), userId).Return(segments, nil).AnyTimes()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/segment/user/", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("userId", fmt.Sprint(userId))
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

		handler.GetSegmentsForUser(w, r)

		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Should return 400 if params is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSegmentSvc := mock_segment.NewMockSegmentService(ctrl)
		handler := NewHandler(nil, mockSegmentSvc)

		var userId = "lol"

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/segment/user/", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("userId", fmt.Sprint(userId))
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

		handler.GetSegmentsForUser(w, r)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Should return 500 if something goes wrong", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var userId int64 = 1

		mockSegmentSvc := mock_segment.NewMockSegmentService(ctrl)
		mockSegmentSvc.EXPECT().GetUserSegments(gomock.Any(), userId).Return(nil, errors.New("internal error"))

		handler := NewHandler(nil, mockSegmentSvc)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/segment/user/", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("userId", fmt.Sprint(userId))
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

		handler.GetSegmentsForUser(w, r)

		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
