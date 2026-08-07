// task パッケージは Task Feature そのもの: /api/tasks を提供するのに必要な
// ものはすべてこのディレクトリに入っている。Laravel側の関心ごとに1ファイル
// (Controller -> handler.go、Requests+Inputs -> request.go、
// Resources -> response.go、UseCases -> usecase.go)という対応関係。
// 新しい Feature を作るときはこのパッケージをコピーして出発点にし、最初の
// 本物の Feature ができたら削除する — Laravel版の app/Features/Task/ と
// まったく同じ使い方。
package task

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/st-man-hori/go-feature-template/internal/httpx"
)

// --- Index -------------------------------------------------------------

// IndexTaskRequest は Requests/IndexTaskRequest.php + Inputs/IndexTaskInput.php に
// 相当する。Laravel の FormRequest は toInput() の中で $this->boolean()/
// $this->integer() を使ってクエリパラメータを読むが、Go には自動的な
// クエリ文字列バインディングが無いため、この template ではたった2フィールドの
// ために reflection ベースのクエリデコーダーを持ち込まず、手で解析している。
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

// StoreTaskRequest は Requests/StoreTaskRequest.php + Inputs/StoreTaskInput.php を
// 1つの型に統合したもの: `json` + `validate` タグを持つ Go の構造体は、
// バリデーション済みのリクエストの形と UseCase への入力の両方をそのまま
// 兼ねられるので、別レイヤーの Input DTO が要らない(Laravel がなぜそれを
// 必要とし Go がなぜ不要かは、このファイル冒頭のパッケージ doc を参照)。
type StoreTaskRequest struct {
	Title       string  `json:"title" validate:"required,max=255"`
	Description *string `json:"description" validate:"omitempty"`
	DueDate     *string `json:"dueDate" validate:"omitempty,datetime=2006-01-02"`
}

// --- Update ----------------------------------------------------------------

// UpdateTaskRequest は Requests/UpdateTaskRequest.php + Inputs/UpdateTaskInput.php に
// 相当する。すべてのフィールドはただのポインタではなく Field[T](field.go)を
// 使っており、これにより PATCH リクエストで「省略された」のか「明示的に null に
// 指定された」のかを区別できる。
type UpdateTaskRequest struct {
	Title       Field[string]  `json:"title"`
	Description Field[*string] `json:"description"`
	DueDate     Field[*string] `json:"dueDate"`
	IsDone      Field[bool]    `json:"isDone"`
}

// Validate は UpdateTaskRequest::rules() 相当のものを手書きしたもの: 実際に
// Set されたフィールドだけをチェックし、Laravel の `sometimes` ルールと同じ
// 挙動になる。go-playground/validator の構造体タグはジェネリクスの Field[T] の
// 内側の Value まできれいには届かないため、この template ではたった1エンドポイント
// のために reflection ベースの validator と格闘するのではなく、Update だけは
// 手動でバリデーションしている。
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
