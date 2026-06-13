package post

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Title     string
	Content   string
	Tags      []string
	CreatedAt time.Time
	UpdatedAt time.Time
}
