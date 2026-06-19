package testutil

import (
	"testing"

	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/platform/db"
	"github.com/stretchr/testify/require"
)

func (e *Env) CreateFollow(t *testing.T, followerID, followeeID uuid.UUID) {
	t.Helper()

	err := e.Q.FollowUser(t.Context(), db.FollowUserParams{
		FollowerID: followerID,
		FolloweeID: followeeID,
	})
	require.NoError(t, err)
}
