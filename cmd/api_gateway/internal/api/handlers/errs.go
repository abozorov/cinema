package handlers

import (
	"errors"
	"net/http"

	"github.com/abozorov/cinema/pkg/errs"
)

// pkg.errs -> http.Error
func errsToHttp(w http.ResponseWriter, err error) {

	switch {
	// http.StatusNotFound
	case errors.Is(err, errs.ErrNotFound):
		http.Error(w, errs.ErrNotFound.Error(), http.StatusNotFound)

	// http.StatusBadRequest
	case errors.Is(err, errs.ErrBadRequest):
		http.Error(w, errs.ErrBadRequest.Error(), http.StatusBadRequest)
	case errors.Is(err, errs.ErrBadRequestBody):
		http.Error(w, errs.ErrBadRequestBody.Error(), http.StatusBadRequest)
	case errors.Is(err, errs.ErrBadRequestQuery):
		http.Error(w, errs.ErrBadRequestQuery.Error(), http.StatusBadRequest)
	case errors.Is(err, errs.ErrVerifyingFailed):
		http.Error(w, errs.ErrVerifyingFailed.Error(), http.StatusBadRequest)
	case errors.Is(err, errs.ErrUserNotBeenVerified):
		http.Error(w, errs.ErrUserNotBeenVerified.Error(), http.StatusBadRequest)
	case errors.Is(err, errs.ErrIncorrectOTPCode):
		http.Error(w, errs.ErrIncorrectOTPCode.Error(), http.StatusBadRequest)

	// http.StatusTooManyRequests
	case errors.Is(err, errs.ErrTooManyRequests):
		http.Error(w, errs.ErrTooManyRequests.Error(), http.StatusTooManyRequests)
	case errors.Is(err, errs.ErrToManyAttempt):
		http.Error(w, errs.ErrToManyAttempt.Error(), http.StatusTooManyRequests)

	// http.StatusGatewayTimeout
	case errors.Is(err, errs.ErrTimeoutExceeded):
		http.Error(w, errs.ErrTimeoutExceeded.Error(), http.StatusGatewayTimeout)

	// http.StatusUnauthorized
	case errors.Is(err, errs.ErrIncorrectLoginOrPassword):
		http.Error(w, errs.ErrIncorrectLoginOrPassword.Error(), http.StatusUnauthorized)
	case errors.Is(err, errs.ErrIncorrectPassword):
		http.Error(w, errs.ErrIncorrectPassword.Error(), http.StatusUnauthorized)
	case errors.Is(err, errs.ErrInvalidToken):
		http.Error(w, errs.ErrInvalidToken.Error(), http.StatusUnauthorized)

	// http.StatusInternalServerError
	default:
		http.Error(w, errs.ErrSomethingWentWrong.Error(), http.StatusInternalServerError)

	}
}
