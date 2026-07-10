package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/abozorov/cinema/pkg/errs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// GORM -> pkg.errs
func postgresToErrs(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, pgx.ErrNoRows): // no rows
		return fmt.Errorf("error %w: %w", err, errs.ErrNotFound)

	case errors.Is(err, context.DeadlineExceeded):
		return fmt.Errorf("error %w: %w", err, errs.ErrTimeoutExceeded)

	case errors.Is(err, context.Canceled):
		return fmt.Errorf("error %w: %w", err, errs.ErrTimeoutExceeded)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return fmt.Errorf("error %w: %w", err, errs.ErrBadRequest)

		case "23503": // foreign_key_violation
			return fmt.Errorf("error %w: %w", err, errs.ErrBadRequest)

		case "23502": // not_null_violation
			return fmt.Errorf("error %w: %w", err, errs.ErrBadRequestBody)

		case "23514": // check_violation
			return fmt.Errorf("error %w: %w", err, errs.ErrBadRequestBody)

		case "22001": // string_data_right_truncation
			return fmt.Errorf("error %w: %w", err, errs.ErrBadRequestBody)

		case "22003": // numeric_value_out_of_range
			return fmt.Errorf("error %w: %w", err, errs.ErrBadRequestBody)

		case "57014": // query_canceled
			return fmt.Errorf("error %w: %w", err, errs.ErrTimeoutExceeded)

		case "22P02": // invalid_text_representation
			return fmt.Errorf("error %w: %w", err, errs.ErrBadRequest)
		}
	}

	return fmt.Errorf("error %w: %w", err, errs.ErrSomethingWentWrong)
}
