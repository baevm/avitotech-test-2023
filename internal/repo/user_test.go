package repo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createUser(t *testing.T, repo *User) int64 {
	userId, err := repo.CreateUser(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, userId)

	return userId
}

func Test_CreateUser(t *testing.T) {
	repo := NewUser(testDbInstance)

	createUser(t, repo)
}

func Test_CheckUserExist(t *testing.T) {
	repo := NewUser(testDbInstance)

	userId, err := repo.CreateUser(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, userId)

	exists, err := repo.CheckUserExist(context.Background(), userId)
	require.NoError(t, err)
	require.True(t, exists)
}
