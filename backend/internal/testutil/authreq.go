package testutil

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/mrbananaaa/gosocialize/internal/platform/requestctx"
)

func AuthenticatedRequest(
	req *http.Request,
	userID uuid.UUID,
) *http.Request {
	ctx := requestctx.WithUser(
		req.Context(),
		requestctx.User{
			ID: userID,
		},
	)

	return req.WithContext(ctx)
}
