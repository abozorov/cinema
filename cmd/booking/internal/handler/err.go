package handler

import (
	"errors"
	"fmt"

	"github.com/abozorov/cinema/cmd/booking/internal/models"
	"github.com/abozorov/cinema/pkg/errs"
	"google.golang.org/grpc/codes"
)

func responseErr(err error) error {
	switch {

	// pkg.errs errors
	case errors.Is(err, errs.ErrBadRequest),
		errors.Is(err, errs.ErrBadRequestBody):
		return fmt.Errorf("error: %d", codes.InvalidArgument)
	case errors.Is(err, errs.ErrTimeoutExceeded):
		return fmt.Errorf("error: %d", codes.DeadlineExceeded)
	case errors.Is(err, errs.ErrNotFound):
		return fmt.Errorf("error: %d", codes.NotFound)

	//model user errors
	case errors.Is(err, models.ErrInvalidBookingUserId),
		errors.Is(err, models.ErrInvalidBookingMovieId),
		errors.Is(err, models.ErrInvalidBookingStatus),
		errors.Is(err, models.ErrEmptyBookingStatus):
		return fmt.Errorf("error: %d", codes.InvalidArgument)

	// default
	default:
		return fmt.Errorf("error: %d", codes.Internal)
	}
}
