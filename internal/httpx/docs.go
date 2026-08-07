package httpx

// ErrorResponse と ValidationErrorResponse は、Error() が実行時に書き込む
// JSON の形を表現するための型。存在意義は、swaggo の @Failure アノテーションが
// 指す先として具象型が必要なことだけで、Handler がこれらを直接構築することは
// ない(Error() はその場で map を組み立てている)。
type ErrorResponse struct {
	Error struct {
		Code    string `json:"code" example:"not_found"`
		Message string `json:"message" example:"Task not found"`
	} `json:"error"`
}

type ValidationErrorResponse struct {
	Message string              `json:"message" example:"The given data was invalid."`
	Errors  map[string][]string `json:"errors"`
}
