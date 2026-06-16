package user

import (
	"github.com/mrbananaaa/gosocialize/store"
)

type Service struct {
	store *store.Store
}

func NewService(
	s *store.Store,
) *Service {
	return &Service{
		store: s,
	}
}
