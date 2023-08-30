package segment

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dezzerlol/avitotech-test-2023/internal/db/models"
	mock_segment "github.com/dezzerlol/avitotech-test-2023/internal/handlers/segment/mocks"
	"github.com/dezzerlol/avitotech-test-2023/internal/repo"
	"github.com/go-chi/chi/v5"
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
}

func Test_DeleteSegment(t *testing.T) {
	t.Run("Should return 200 and delete segment", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSegmentSvc := mock_segment.NewMockSegmentService(ctrl)
		handler := NewHandler(nil, mockSegmentSvc)

		segment := &models.Segment{
			Slug: "TEST_SEGMENT",
		}

		mockSegmentSvc.EXPECT().DeleteBySlug(gomock.Any(), segment).Return(nil).AnyTimes()

		body := `{"slug": "TEST_SEGMENT"}`

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/segment", strings.NewReader(body))
		handler.Delete(w, r)

		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Should return 404 if segment not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSegmentSvc := mock_segment.NewMockSegmentService(ctrl)
		handler := NewHandler(nil, mockSegmentSvc)

		segment := &models.Segment{
			Slug: "TEST_SEGMENT",
		}

		mockSegmentSvc.EXPECT().DeleteBySlug(gomock.Any(), segment).Return(repo.ErrSegmentNotFound)

		body := `{"slug": "TEST_SEGMENT"}`

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/segment", strings.NewReader(body))
		handler.Delete(w, r)

		require.Equal(t, http.StatusNotFound, w.Code)
	})
}

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
		rctx.URLParams.Add("userId", "1")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

		handler.GetSegmentsForUser(w, r)

		require.Equal(t, http.StatusOK, w.Code)
	})
}

func Test_UpdateUserSegments(t *testing.T) {
	t.Run("Should return 200 and add user segments", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSegmentSvc := mock_segment.NewMockSegmentService(ctrl)
		handler := NewHandler(nil, mockSegmentSvc)

		addSegments := []string{"TEST_SEGMENT1", "TEST_SEGMENT2", "TEST_SEGMENT3"}
		var userId int64 = 1

		mockSegmentSvc.EXPECT().
			UpdateUserSegments(gomock.Any(), userId, addSegments, int64(0), nil).
			Return(int64(len(addSegments)), int64(0), nil).AnyTimes()

		body := UpdateUserSegmentsRequest{
			AddSegments: addSegments,
			UserId:      1,
			TTL:         0,
		}

		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/segment/user/", &buf)

		handler.UpdateUserSegments(w, r)

		require.Equal(t, http.StatusOK, w.Code)
		require.JSONEq(t, `{"segments_added": 3, "segments_deleted": 0}`, w.Body.String())
	})
}
