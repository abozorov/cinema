package handler

import (
	"errors"

	"github.com/abozorov/cinema/cmd/user/internal/models"
	"github.com/abozorov/cinema/pkg/errs"
)

func responseErr(err error) error {
	switch {

	// pkg.errs errors
	case errors.Is(err, errs.ErrBadRequest):
		return errs.ErrBadRequest
	case errors.Is(err, errs.ErrBadRequestBody):
		return errs.ErrBadRequestBody
	case errors.Is(err, errs.ErrTimeoutExceeded):
		return errs.ErrTimeoutExceeded
	case errors.Is(err, errs.ErrNotFound):
		return errs.ErrNotFound

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
		return errs.ErrBadRequest

	// default
	default:
		return errs.ErrSomethingWentWrong
	}
}
