package auth_test

import (
	"testing"

	"github.com/mrbananaaa/gosocialize/internal/auth"
	"github.com/mrbananaaa/gosocialize/internal/platform/jwt"
	"github.com/mrbananaaa/gosocialize/internal/testutil"
	"github.com/mrbananaaa/gosocialize/internal/user"
	"github.com/mrbananaaa/gosocialize/store"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	_, svc := newAuthService(t)

	t.Run("success - user registered", func(t *testing.T) {
		u := user.User{
			Email:    "mailme@mrbananaaid.com",
			Username: "mrbananaaaa",
			Password: "supersecurepassword",
			Name:     "Mr Bananaaa",
		}

		user, err := svc.Register(t.Context(), auth.RegisterInput{
			Email:    u.Email,
			Username: u.Username,
			Password: u.Password,
			Name:     u.Name,
		})
		require.NoError(t, err)

		require.Equal(t, u.Email, user.Email)
		require.Equal(t, u.Username, user.Username)
		require.Equal(t, u.Name, user.Name)
		require.NotEmpty(t, user.ID)
		require.NotEmpty(t, user.Password)
		require.NotEmpty(t, user.CreatedAt)
		require.NotEmpty(t, user.UpdatedAt)
	})
}

func newAuthService(t *testing.T) (*testutil.Env, *auth.Service) {
	t.Helper()

	env := testutil.NewEnv(t)
	store := store.New(env.DB)
	tokenSvc := jwt.New()
	hasher := auth.NewArgon2Hasher()
	svc := auth.NewService(tokenSvc, hasher, store)

	return env, svc
}
