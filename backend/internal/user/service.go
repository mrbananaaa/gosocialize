package user

import (
	"github.com/mrbananaaa/gosocialize/internal/platform/database/postgres/sqlc"
)

type Service struct {
	q *sqlc.Queries
}

func NewService(q *sqlc.Queries) *Service {
	return &Service{
		q: q,
	}
}
