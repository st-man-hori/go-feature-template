// Package apperror defines the domain-error type UseCases return for
// business-rule failures (not found, conflict, etc).
//
// Laravel comparison: this plays the role of App\Shared\Exceptions\AppDomainException.
// Laravel uses a class hierarchy (subclass per error, caught by bootstrap/app.php's
// withExceptions renderer). Go has no exceptions, so UseCases return errors as
// values instead of throwing, and *AppError is a single concrete type carrying
// the code/message/status that would otherwise live on each subclass.
package apperror

import "net/http"

// AppError is a domain/business-rule error that the HTTP layer knows how to
// render as a JSON error response (see internal/httpx.WriteError).
type AppError struct {
	Code       string
	Message    string
	HTTPStatus int
}

func (e *AppError) Error() string {
	return e.Message
}

func New(code, message string, httpStatus int) *AppError {
	return &AppError{Code: code, Message: message, HTTPStatus: httpStatus}
}

// NotFound builds the common "record not found" case, mirroring what
// Eloquent's findOrFail() raises automatically in the Laravel template.
func NotFound(resource string) *AppError {
	return New("not_found", resource+" not found", http.StatusNotFound)
}
