package segment

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dezzerlol/avitotech-test-2023/internal/db/models"
	mock_segment "github.com/dezzerlol/avitotech-test-2023/internal/handlers/segment/mocks"
	"github.com/dezzerlol/avitotech-test-2023/internal/repo"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

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

		body := fmt.Sprintf(`{"slug": "%s"}`, segment.Slug)

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

		body := fmt.Sprintf(`{"slug": "%s"}`, segment.Slug)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/segment", strings.NewReader(body))
		handler.Delete(w, r)

		require.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Should return 500 if something goes wrong", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		segment := &models.Segment{
			Slug: "TEST_SEGMENT",
		}

		mockSegmentSvc := mock_segment.NewMockSegmentService(ctrl)
		mockSegmentSvc.EXPECT().DeleteBySlug(gomock.Any(), gomock.Any()).Return(errors.New("internal error"))

		handler := NewHandler(nil, mockSegmentSvc)

		body := fmt.Sprintf(`{"slug": "%s"}`, segment.Slug)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/segment", strings.NewReader(body))
		handler.Delete(w, r)

		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Should return 500 if slug is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSegmentSvc := mock_segment.NewMockSegmentService(ctrl)

		handler := NewHandler(nil, mockSegmentSvc)

		body := `{"slug": "T"}`

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/segment", strings.NewReader(body))
		handler.Delete(w, r)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Should return 400 if JSON is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSegmentSvc := mock_segment.NewMockSegmentService(ctrl)

		handler := NewHandler(nil, mockSegmentSvc)

		body := `{"slug: "TEST_SEGMENT"}`

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/segment", strings.NewReader(body))
		handler.Delete(w, r)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}
