// Package httpx は、各 Feature の Handler が共通して使う小さな HTTP ヘルパーを
// まとめたもの: JSON のエンコード/デコード、{"data": ...} エンベロープ、
// エラーレンダリング。
//
// Laravel比較: request.go と合わせて、JsonResource の response()、
// Illuminate\Http\JsonResponse、そして bootstrap/app.php の withExceptions() で
// 組み込まれる App\Shared\Responses\ApiErrorResponse + AppDomainException の
// レンダリングを合わせた役割を果たす。フレームワークのフックではなくただの
// パッケージなのは、Go には Laravel の例外レンダラーパイプラインに相当するものが
// 無く、Handler が明示的に httpx.Error(w, err) を呼ぶ形になるため。
package httpx

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/st-man-hori/go-feature-template/internal/apperror"
)

// JSON は v を、指定したステータスコードで JSON ボディとして書き込む。
func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(v)
}

// Data は v を {"data": v} でラップする。Laravel の JsonResource や
// AnonymousResourceCollection が自動的に作るエンベロープと同じ形。
func Data(w http.ResponseWriter, status int, v any) {
	JSON(w, status, map[string]any{"data": v})
}

// NoContent は空の 204 を書き込む。TaskController::destroy() と同じ挙動。
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Error は err を JSON エラーレスポンスとしてレンダリングする。ステータスコードは
// err の型によって決まる:
//   - *apperror.AppError  -> 自身が持つ HTTPStatus(例: NotFound なら 404)
//   - *ValidationError    -> 422。Laravel のデフォルトのバリデーションエラー形式
//   - それ以外            -> 500。実際のエラーはログに出すがレスポンスには出さない
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
