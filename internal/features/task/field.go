package task

import "encoding/json"

// Field は、PATCH リクエストにおいて「リクエストボディに存在しない」
// (Set == false)と「明示的に指定された」(Set == true。*T で null の場合もある)
// を区別できる値を表す。encoding/json はボディに実際に存在するキーに対して
// しか UnmarshalJSON を呼ばないため、キーが省略されていれば Set は false のままになる。
//
// Laravel比較: これは spatie/laravel-data の Optional が UpdateTaskInput.php で
// 解決しているのとまったく同じ問題を解く — ただのポインタでは「フィールドが
// 省略された」のか「フィールドが明示的に null にされた」のかを区別できない。
// JSON の `null` とキーの欠落は、どちらも nil ポインタにデコードされてしまうため。
type Field[T any] struct {
	Value T
	Set   bool
}

func (f *Field[T]) UnmarshalJSON(data []byte) error {
	f.Set = true
	return json.Unmarshal(data, &f.Value)
}
