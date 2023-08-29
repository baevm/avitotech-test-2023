package repo

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/dezzerlol/avitotech-test-2023/internal/db/models"
	"github.com/dezzerlol/avitotech-test-2023/pkg/testhelper"
	"github.com/stretchr/testify/require"
)

func createSegment(t *testing.T, repo *Segment) *models.Segment {
	var segment = &models.Segment{
		Slug: testhelper.RandomString(12),
	}

	err := repo.Create(context.Background(), segment)
	require.NoError(t, err)
	require.NotEmpty(t, segment.CreatedAt)
	require.WithinDuration(t, time.Now(), *segment.CreatedAt, 1*time.Second)

	return segment
}

func addUserSegments(t *testing.T, repo *Segment, userId int64) []string {
	var segments = make([]string, 3)

	for i := range segments {
		segment := createSegment(t, repo)
		segments[i] = segment.Slug
	}

	addedSegments, err := repo.AddUserSegments(context.Background(), userId, segments, 0)
	require.NoError(t, err)
	require.Equal(t, len(segments), int(addedSegments))

	return segments
}

func Test_CreateSegment(t *testing.T) {
	repo := NewSegment(testDbInstance)

	createSegment(t, repo)
}

func Test_DeleteSegment(t *testing.T) {
	repo := NewSegment(testDbInstance)

	segment := createSegment(t, repo)

	err := repo.DeleteBySlug(context.Background(), segment)
	require.NoError(t, err)
}

func Test_AddUserSegments(t *testing.T) {
	repo := NewSegment(testDbInstance)

	userId := createUser(t, NewUser(testDbInstance))

	addUserSegments(t, repo, userId)
}

func Test_GetUserSegments(t *testing.T) {
	repo := NewSegment(testDbInstance)

	userId := createUser(t, NewUser(testDbInstance))
	segments := addUserSegments(t, repo, userId)

	userSegments, err := repo.GetUserSegments(context.Background(), userId)

	require.NoError(t, err)
	require.NotEmpty(t, userSegments)

	sort.Strings(segments)

	for i, segment := range userSegments {
		require.Equal(t, segment.Slug, segments[i])
	}
}

func Test_DeleteUserSegments(t *testing.T) {
	repo := NewSegment(testDbInstance)

	userId := createUser(t, NewUser(testDbInstance))
	segments := addUserSegments(t, repo, userId)

	deletedSegments, err := repo.DeleteUserSegments(context.Background(), userId, segments)
	require.NoError(t, err)
	require.Equal(t, len(segments), int(deletedSegments))
}

func Test_GetUserHistory(t *testing.T) {
	repo := NewSegment(testDbInstance)

	userId := createUser(t, NewUser(testDbInstance))
	addUserSegments(t, repo, userId)

	date := time.Now()

	history, err := repo.GetUserHistory(context.Background(), userId, date)
	require.NoError(t, err)
	require.NotEmpty(t, history)
}
