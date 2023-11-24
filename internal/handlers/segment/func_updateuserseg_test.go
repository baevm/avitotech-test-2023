package segment

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mock_segment "github.com/dezzerlol/avitotech-test-2023/internal/handlers/segment/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_UpdateUserSegments(t *testing.T) {
	t.Run("Should return 200 and add user segments", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSegmentSvc := mock_segment.NewMockSegmentService(ctrl)
		handler := NewHandler(nil, mockSegmentSvc)

		addSegments := []string{"TEST_SEGMENT1", "TEST_SEGMENT2", "TEST_SEGMENT3"}
		var userId int64 = 1
		var ttl int64 = 0

		body := UpdateUserSegmentsRequest{
			AddSegments: addSegments,
			UserId:      userId,
			TTL:         ttl,
		}

		mockSegmentSvc.EXPECT().
			UpdateUserSegments(gomock.Any(), userId, addSegments, ttl, nil).
			Return(int64(len(addSegments)), int64(0), nil).AnyTimes()

		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/segment/user/", &buf)

		handler.UpdateUserSegments(w, r)

		require.Equal(t, http.StatusOK, w.Code)
		require.JSONEq(t, `{"segments_added": 3, "segments_deleted": 0}`, w.Body.String())
	})

	t.Run("Should return 400 if JSON is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSegmentSvc := mock_segment.NewMockSegmentService(ctrl)

		handler := NewHandler(nil, mockSegmentSvc)

		body := `{user_id: 1, "ttl": 1, "add_segments": ["TEST_SEGMENT1", "TEST_SEGMENT2", "TEST_SEGMENT3"]}`

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/segment/user", strings.NewReader(body))
		handler.UpdateUserSegments(w, r)

		fmt.Println(w.Body)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Should return 400 if body is invalid after validation", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSegmentSvc := mock_segment.NewMockSegmentService(ctrl)

		addSegments := []string{"TEST_SEGMENT1", "TEST_SEGMENT2", "TEST_SEGMENT3"}
		var userId int64 = -100
		var ttl int64 = 0

		body := UpdateUserSegmentsRequest{
			AddSegments: addSegments,
			UserId:      userId,
			TTL:         ttl,
		}

		handler := NewHandler(nil, mockSegmentSvc)

		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/segment/user", &buf)
		handler.UpdateUserSegments(w, r)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Should return 500 if something goes wrong", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		addSegments := []string{"TEST_SEGMENT1", "TEST_SEGMENT2", "TEST_SEGMENT3"}
		var userId int64 = 1
		var ttl int64 = 0

		body := UpdateUserSegmentsRequest{
			AddSegments: addSegments,
			UserId:      userId,
			TTL:         ttl,
		}

		mockSegmentSvc := mock_segment.NewMockSegmentService(ctrl)
		mockSegmentSvc.EXPECT().
			UpdateUserSegments(gomock.Any(), userId, addSegments, ttl, nil).
			Return(int64(0), int64(0), errors.New("internal error")).AnyTimes()

		handler := NewHandler(nil, mockSegmentSvc)

		var buf bytes.Buffer
		json.NewEncoder(&buf).Encode(body)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/segment/user", &buf)
		handler.UpdateUserSegments(w, r)

		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
