// Package httpx holds small HTTP helpers shared by every feature's handler:
// JSON encoding/decoding, the {"data": ...} envelope, and error rendering.
//
// Laravel comparison: together with request.go, this package plays the role
// of JsonResource's response(), Illuminate\Http\JsonResponse, and the
// App\Shared\Responses\ApiErrorResponse + AppDomainException rendering wired
// up in bootstrap/app.php's withExceptions(). It's a package, not a
// framework hook, because Go has no equivalent to Laravel's exception
// renderer pipeline — handlers call httpx.Error(w, err) explicitly instead.
package httpx

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/st-man-hori/go-feature-template/internal/apperror"
)

// JSON writes v as a JSON body with the given status code.
func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(v)
}

// Data wraps v in {"data": v}, mirroring the envelope Laravel's JsonResource
// and AnonymousResourceCollection produce automatically.
func Data(w http.ResponseWriter, status int, v any) {
	JSON(w, status, map[string]any{"data": v})
}

// NoContent writes an empty 204, mirroring TaskController::destroy().
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Error renders err as a JSON error response. The status code depends on
// the error's type:
//   - *apperror.AppError  -> its own HTTPStatus (e.g. 404 for NotFound)
//   - *ValidationError    -> 422, in Laravel's default validation-error shape
//   - anything else       -> 500, with the real error logged but not exposed
func Error(w http.ResponseWriter, err error) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		JSON(w, appErr.HTTPStatus, map[string]any{
			"error": map[string]any{
				"code":    appErr.Code,
				"message": appErr.Message,
			},
		})
		return
	}

	var verr *ValidationError
	if errors.As(err, &verr) {
		JSON(w, http.StatusUnprocessableEntity, map[string]any{
			"message": "The given data was invalid.",
			"errors":  verr.Errors,
		})
		return
	}

	slog.Error("unhandled error", "error", err)
	JSON(w, http.StatusInternalServerError, map[string]any{
		"error": map[string]any{
			"code":    "internal_error",
			"message": "Something went wrong.",
		},
	})
}
