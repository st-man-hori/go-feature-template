// Package task is the Task feature: everything needed to serve /api/tasks
// lives in this directory, one file per Laravel-side concern (Controller ->
// handler.go, Requests+Inputs -> request.go, Resources -> response.go,
// UseCases -> usecase.go). Copy this package as a starting point for a new
// Feature, then delete it once your first real Feature exists — exactly
// like the Laravel template's app/Features/Task/.
package task

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/st-man-hori/go-feature-template/internal/httpx"
)

// --- Index -------------------------------------------------------------

// IndexTaskRequest mirrors Requests/IndexTaskRequest.php + Inputs/IndexTaskInput.php.
// Laravel's FormRequest reads query params via $this->boolean()/$this->integer()
// inside toInput(); Go has no automatic query-string binding, so this
// template parses the two fields by hand rather than pulling in a
// reflection-based query decoder for a single endpoint.
type IndexTaskRequest struct {
	IsDone  *bool
	PerPage int
}

func NewIndexTaskRequest(query url.Values) (IndexTaskRequest, error) {
	req := IndexTaskRequest{PerPage: 15}

	if raw := query.Get("isDone"); raw != "" {
		v, err := strconv.ParseBool(raw)
		if err != nil {
			return req, &httpx.ValidationError{Errors: map[string][]string{"isDone": {"This field must be a boolean."}}}
		}
		req.IsDone = &v
	}

	if raw := query.Get("perPage"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			return req, &httpx.ValidationError{Errors: map[string][]string{"perPage": {"This field must be an integer."}}}
		}
		req.PerPage = v
	}

	if req.PerPage < 1 || req.PerPage > 100 {
		return req, &httpx.ValidationError{Errors: map[string][]string{"perPage": {"This field must be between 1 and 100."}}}
	}

	return req, nil
}

// --- Store ---------------------------------------------------------------

// StoreTaskRequest mirrors Requests/StoreTaskRequest.php + Inputs/StoreTaskInput.php
// merged into one type: a Go struct with `json` + `validate` tags already
// works as both the validated request shape and the UseCase input, so there
// is no separate Input DTO layer (see request.go's package doc for why
// Laravel needs one and Go doesn't).
type StoreTaskRequest struct {
	Title       string  `json:"title" validate:"required,max=255"`
	Description *string `json:"description" validate:"omitempty"`
	DueDate     *string `json:"dueDate" validate:"omitempty,datetime=2006-01-02"`
}

// --- Update ----------------------------------------------------------------

// UpdateTaskRequest mirrors Requests/UpdateTaskRequest.php + Inputs/UpdateTaskInput.php.
// Every field uses Field[T] (field.go) instead of a plain pointer so a PATCH
// request can tell "omitted" apart from "explicitly set to null".
type UpdateTaskRequest struct {
	Title       Field[string]  `json:"title"`
	Description Field[*string] `json:"description"`
	DueDate     Field[*string] `json:"dueDate"`
	IsDone      Field[bool]    `json:"isDone"`
}

// Validate hand-rolls the equivalent of UpdateTaskRequest::rules(): only
// fields that were actually Set get checked, matching Laravel's `sometimes`
// rule. go-playground/validator's struct tags don't reach into a generic
// Field[T]'s inner Value cleanly, so this template validates Update by hand
// rather than fighting the reflection-based validator for one endpoint.
func (r UpdateTaskRequest) Validate() error {
	errs := map[string][]string{}

	if r.Title.Set {
		if strings.TrimSpace(r.Title.Value) == "" {
			errs["title"] = append(errs["title"], "This field is required.")
		} else if len(r.Title.Value) > 255 {
			errs["title"] = append(errs["title"], "This field must be at most 255 characters.")
		}
	}

	if r.DueDate.Set && r.DueDate.Value != nil {
		if _, err := time.Parse("2006-01-02", *r.DueDate.Value); err != nil {
			errs["dueDate"] = append(errs["dueDate"], "This field must be a valid date (YYYY-MM-DD).")
		}
	}

	if len(errs) > 0 {
		return &httpx.ValidationError{Errors: errs}
	}
	return nil
}
