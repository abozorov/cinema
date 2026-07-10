package errs

import "errors"

var (
	// ServerError
	ErrSomethingWentWrong = errors.New("something went wrong")

	// BadRequest
	ErrBadRequest      = errors.New("bad request")
	ErrBadRequestBody  = errors.New("bad request body")
	ErrBadRequestQuery = errors.New("bad request query")

	// Too Many Requests
	ErrTooManyRequests = errors.New("too many requests")

	// Timeout exceeded
	ErrTimeoutExceeded = errors.New("timeout exceeded")

	// Verify
	ErrVerifyingFailed     = errors.New("Verifying failed")
	ErrUserNotBeenVerified = errors.New("email has not been verified yet")
	ErrIncorrectOTPCode    = errors.New("incorrect OTP code")
	ErrToManyAttempt       = errors.New("user is too many attempts")

	// Unautorized
	ErrIncorrectLoginOrPassword = errors.New("incorrect login or password")
	ErrIncorrectPassword        = errors.New("incorrect password")
	ErrInvalidToken             = errors.New("token is invalid")

	// NotFound
	ErrNotFound = errors.New("not found")
)
