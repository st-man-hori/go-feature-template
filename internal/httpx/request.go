package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = newValidator()

func newValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())

	// Report validation errors using each field's `json` tag (e.g. "dueDate")
	// instead of its Go field name ("DueDate"), so error responses match the
	// camelCase the API accepts and returns.
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return v
}

// ValidationError is returned by Decode when the request body fails
// validation. httpx.Error renders it as Laravel's default 422 shape:
// {"message": ..., "errors": {field: [msg, ...]}}.
type ValidationError struct {
	Errors map[string][]string
}

func (e *ValidationError) Error() string {
	return "validation failed"
}

// Decode reads a JSON body into dst, then validates it against dst's
// `validate` struct tags.
//
// Laravel comparison: this collapses FormRequest::rules() +
// Data::from($request->validated()) into one step. A Laravel FormRequest is
// split from its Input DTO because the FormRequest is tied to
// Illuminate\Http\Request (it needs authorize(), $this->boolean(), etc.)
// while the Input DTO is plain data passed to the UseCase. In Go, a plain
// struct with `json` + `validate` tags already serves both roles, so each
// Task request type in this template plays both parts at once.
func Decode(r *http.Request, dst any) error {
	if err := DecodeJSON(r, dst); err != nil {
		return err
	}
	return Validate(dst)
}

// DecodeJSON reads a JSON body into dst without running struct-tag
// validation, for requests (like Update) that validate by hand instead —
// see UpdateTaskRequest.Validate for why.
func DecodeJSON(r *http.Request, dst any) error {
	if r.Body == nil {
		return &ValidationError{Errors: map[string][]string{"body": {"request body is required"}}}
	}
	defer func() { _ = r.Body.Close() }()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return &ValidationError{Errors: map[string][]string{"body": {"request body must be valid JSON"}}}
	}

	return nil
}

// Validate runs struct-tag validation without decoding a body first, for
// requests (like Index) whose input comes from query parameters instead.
func Validate(dst any) error {
	if err := validate.Struct(dst); err != nil {
		return toValidationError(err)
	}
	return nil
}

func toValidationError(err error) *ValidationError {
	fieldErrors := map[string][]string{}

	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) {
		for _, fe := range verrs {
			fieldErrors[fe.Field()] = append(fieldErrors[fe.Field()], validationMessage(fe))
		}
	}

	return &ValidationError{Errors: fieldErrors}
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required."
	case "max":
		return fmt.Sprintf("This field must be at most %s characters.", fe.Param())
	case "min":
		return fmt.Sprintf("This field must be at least %s.", fe.Param())
	case "datetime":
		return "This field must be a valid date."
	default:
		return fmt.Sprintf("This field failed validation on the '%s' rule.", fe.Tag())
	}
}
