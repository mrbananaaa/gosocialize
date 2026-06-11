package pagination

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type CursorQueryParam struct {
	Cursor string
	Limit  int32
}

type PaginationMeta struct {
	NextCursor string `json:"next_cursor"`
	HasMore    bool   `json:"has_more"`
}

type Cursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        uuid.UUID `json:"id"`
}

func EncodeCursor(c Cursor) (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	encodedCursor := base64.StdEncoding.EncodeToString(b)

	return encodedCursor, nil
}

func DecodeCursor(raw string) (*Cursor, error) {
	if raw == "" {
		return nil, nil
	}

	b, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, err
	}

	var c Cursor
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
