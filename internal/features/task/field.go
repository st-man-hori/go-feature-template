package task

import "encoding/json"

// Field represents a PATCH-friendly value that distinguishes "not present in
// the request body" (Set == false) from "explicitly provided" (Set == true,
// possibly to null via *T). encoding/json only calls UnmarshalJSON for keys
// that actually appear in the body, so an omitted key leaves Set false.
//
// Laravel comparison: this solves the exact problem spatie/laravel-data's
// Optional solves in UpdateTaskInput.php — a plain pointer can't tell
// "field omitted" apart from "field explicitly set to null", since JSON
// `null` and an absent key both decode to a nil pointer.
type Field[T any] struct {
	Value T
	Set   bool
}

func (f *Field[T]) UnmarshalJSON(data []byte) error {
	f.Set = true
	return json.Unmarshal(data, &f.Value)
}
