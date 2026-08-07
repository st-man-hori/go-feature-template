// Package apperror は、UseCase がビジネスルール違反(not found, conflict など)
// を返すためのドメインエラー型を定義する。
//
// Laravel比較: App\Shared\Exceptions\AppDomainException の役割を果たす。
// Laravel はクラス階層(エラーごとにサブクラスを作り、bootstrap/app.php の
// withExceptions レンダラーで捕捉する)を使うが、Go には例外機構が無いため、
// UseCase は throw する代わりに値として error を返す。*AppError は、各サブクラス
// に分散していたはずの code/message/status を1つの具象型にまとめたもの。
package apperror

import "net/http"

// AppError はドメイン/ビジネスルールのエラーで、HTTP 層(internal/httpx.Error)が
// JSON エラーレスポンスとしてレンダリングする方法を知っている型。
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

// NotFound は「レコードが見つからない」という頻出パターンを組み立てる。
// Laravel版で Eloquent の findOrFail() が自動的に投げるものと同じ役割。
func NotFound(resource string) *AppError {
	return New("not_found", resource+" not found", http.StatusNotFound)
}
