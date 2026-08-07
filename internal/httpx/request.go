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

	// バリデーションエラーは Go のフィールド名("DueDate")ではなく、各フィールドの
	// `json` タグ(例: "dueDate")を使って報告する。これにより、エラーレスポンスが
	// API が受け取り・返す camelCase と一致する。
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return v
}

// ValidationError は、リクエストボディがバリデーションに失敗したときに Decode が
// 返す。httpx.Error はこれを Laravel のデフォルトの 422 形式
// {"message": ..., "errors": {field: [msg, ...]}} でレンダリングする。
type ValidationError struct {
	Errors map[string][]string
}

func (e *ValidationError) Error() string {
	return "validation failed"
}

// Decode は JSON ボディを dst に読み込み、dst の `validate` 構造体タグに従って
// バリデーションする。
//
// Laravel比較: これは FormRequest::rules() + Data::from($request->validated())
// を1ステップにまとめたもの。Laravel の FormRequest が Input DTO と分かれて
// いるのは、FormRequest が Illuminate\Http\Request に紐づいている
// (authorize() や $this->boolean() などが必要)一方、Input DTO は UseCase に
// 渡すただのデータだからである。Go では `json` + `validate` タグを持つただの
// 構造体がその両方の役割を兼ねられるため、この template の各 Task リクエスト型は
// 両方の役割を1つで担っている。
func Decode(r *http.Request, dst any) error {
	if err := DecodeJSON(r, dst); err != nil {
		return err
	}
	return Validate(dst)
}

// DecodeJSON は構造体タグによるバリデーションを行わずに JSON ボディを dst に
// 読み込む。Update のように手動でバリデーションするリクエスト向け
// — 理由は UpdateTaskRequest.Validate を参照。
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

// Validate はボディのデコードを行わず、構造体タグによるバリデーションだけを
// 実行する。Index のように、入力がクエリパラメータから来るリクエスト向け。
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
