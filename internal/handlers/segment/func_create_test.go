package segment

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dezzerlol/avitotech-test-2023/internal/db/models"
	mock_segment "github.com/dezzerlol/avitotech-test-2023/internal/handlers/segment/mocks"
	"github.com/dezzerlol/avitotech-test-2023/internal/repo"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_CreateSegment(t *testing.T) {
	t.Run("Should return 201 and create segment", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		segment := &models.Segment{
			Slug: "TEST_SEGMENT",
		}

		mockSegmentSvc := mock_segment.NewMockSegmentService(ctrl)

		mockSegmentSvc.EXPECT().Create(gomock.Any(), segment).DoAndReturn(func(ctx context.Context, segment *models.Segment) error {
			segment.CreatedAt = &time.Time{}
			return nil
		})

		handler := NewHandler(nil, mockSegmentSvc)

		body := `{"slug": "TEST_SEGMENT"}`

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/segment", strings.NewReader(body))
		handler.Create(w, r)

		require.Equal(t, http.StatusCreated, w.Code)

		createdAtStr := w.Body.String()

		createdAt, _ := time.Parse(time.RFC3339, createdAtStr)
		require.WithinDuration(t, time.Time{}, createdAt, 5*time.Second)
	})

	t.Run("Should return 400 if slug is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSegmentSvc := mock_segment.NewMockSegmentService(ctrl)

		mockSegmentSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

		handler := NewHandler(nil, mockSegmentSvc)

		body := `{"slug": "T"}`

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/segment", strings.NewReader(body))
		handler.Create(w, r)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Should return 400 if slug already exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSegmentSvc := mock_segment.NewMockSegmentService(ctrl)
		mockSegmentSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(repo.ErrSegmentAlreadyExists)

		handler := NewHandler(nil, mockSegmentSvc)

		body := `{"slug": "TEST_SEGMENT"}`

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/segment", strings.NewReader(body))
		handler.Create(w, r)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Should return 500 if something goes wrong", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSegmentSvc := mock_segment.NewMockSegmentService(ctrl)
		mockSegmentSvc.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("internal error"))

		handler := NewHandler(nil, mockSegmentSvc)

		body := `{"slug": "TEST_SEGMENT"}`

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/segment", strings.NewReader(body))
		handler.Create(w, r)

		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Should return 400 if JSON is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSegmentSvc := mock_segment.NewMockSegmentService(ctrl)

		handler := NewHandler(nil, mockSegmentSvc)

		body := `{"slug: "TEST_SEGMENT"}`

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/segment", strings.NewReader(body))
		handler.Create(w, r)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}
