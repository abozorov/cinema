package handler

import (
	"errors"
	"fmt"

	"github.com/abozorov/cinema/cmd/user/internal/models"
	"github.com/abozorov/cinema/pkg/errs"
	"google.golang.org/grpc/codes"
)

func responseErr(err error) error {
	switch {

	// pkg.errs errors
	case errors.Is(err, errs.ErrBadRequest):
		return fmt.Errorf("error: %d", codes.InvalidArgument)
	case errors.Is(err, errs.ErrBadRequestBody):
		return fmt.Errorf("error: %d", codes.InvalidArgument)
	case errors.Is(err, errs.ErrTimeoutExceeded):
		return fmt.Errorf("error: %d", codes.DeadlineExceeded)
	case errors.Is(err, errs.ErrNotFound):
		return fmt.Errorf("error: %d", codes.NotFound)

	//model user errors
	case errors.Is(err, models.ErrEmptyName),
		errors.Is(err, models.ErrEmptyEmail),
		errors.Is(err, models.ErrEmptyPhone),
		errors.Is(err, models.ErrEmptyUserID),
		errors.Is(err, models.ErrInvalidAge),
		errors.Is(err, models.ErrInvalidEmail),
		errors.Is(err, models.ErrInvalidName),
		errors.Is(err, models.ErrInvalidPhone),
		errors.Is(err, models.ErrInvalidUSerId):
		return fmt.Errorf("error: %d", codes.InvalidArgument)

	// default
	default:
		return fmt.Errorf("error: %d", codes.Internal)
	}
}
