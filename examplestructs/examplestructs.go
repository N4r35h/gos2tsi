package examplestructs

import "time"

type SimpleStruct1 struct {
	StringField    string  `json:"string_field"`
	BooleanField   bool    `json:"boolean_field"`
	UINTField      uint    `json:"uint_field"`
	INTField       int     `json:"int_field"`
	Float64Field   float64 `json:"float64_field"`
	FieldWOJSONTag string
	HiddenField    string `json:"-"`
}

type EmbedableStruct struct {
	ID uint `json:"id"`
}

type StructWithEmbeding struct {
	EmbedableStruct
	Level string `json:"level"`
}

type SimpleStruct struct {
	Test string `json:"test"`
}

type StructWithFieldStruct struct {
	StructField SimpleStruct `json:"struct_field"`
}

type StructWithFieldSlice struct {
	SliceField []string `json:"slice_field"`
}

type StructWithFieldStructSlice struct {
	StructSliceField []SimpleStruct `json:"struct_slice_field"`
}

type SingleGenericStruct[T any] struct {
	Data           []T          `json:"data"`
	PaginationInfo SimpleStruct `json:"pagination_info"`
}

type MultiGenericStruct[T any, U any] struct {
	Data           []T          `json:"data"`
	Data2          []U          `json:"data2"`
	PaginationInfo SimpleStruct `json:"pagination_info"`
}

// Test
// Test 123
/* Test 123 */
type StructWithCommentOnTop struct {
	Test string `json:"test"`
}

type StructWithCustomTSType struct {
	CustomeTSType time.Time `json:"custom_ts_type" ts_type:"string"`
}
