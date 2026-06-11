package postgres

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mrbananaaa/gosocialize/internal/platform/apperr"
)

// 23505	unique violation
// 23503	foreign key violation
// 23502	not null violation
// 23514	check constraint violation
// 22P02	invalid input syntax
// 40001	serialization failure
// 40P01	deadlock detected

const (
	UniqueViolation      = "23505"
	ForeignKeyViolation  = "23503"
	NotNullViolation     = "23502"
	ConstraintViolation  = "23514"
	InvalidInput         = "22P02"
	SerializationFailure = "40001"
	Deadlock             = "40P01"
)

func PgErrMapper(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return apperr.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case ForeignKeyViolation:
			return apperr.Wrap(
				pgErr,
				apperr.Code.UserNotFound,
				fmt.Sprintf("%s are not found", resolveUniqueField(pgErr)),
			)

		case UniqueViolation:
			constraint := resolveUniqueField(pgErr)

			return apperr.Wrap(
				pgErr,
				apperr.Code.Conflict,
				fmt.Sprintf("%s already exists", constraint),
			)
		}
	}

	// fallback error
	return apperr.ErrInternal
}

func resolveUniqueField(pgErr *pgconn.PgError) string {
	name := pgErr.ConstraintName

	parts := strings.Split(name, "_")
	if len(parts) >= 2 {
		return parts[len(parts)-2]
	}

	return ""
}
