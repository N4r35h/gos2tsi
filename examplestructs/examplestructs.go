package examplestructs

import (
	"time"

	"github.com/N4r35h/gos2tsi/exstructpkg2"
	"github.com/N4r35h/gos2tsi/exstructpkg3"
)

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

type PrimitiveStruct struct {
	// Primitive types
	Boolean   bool        `json:"boolean"`
	Interface interface{} `json:"interface"`
	String    string      `json:"string"`
	Random    any         `json:"random"`
	Int       int         `json:"int"`
	Int8      int8        `json:"int8"`
	Int16     int16       `json:"int16"`
	Int32     int32       `json:"int32"`
	Int64     int64       `json:"int64"`
	Uint      uint        `json:"uint"`
	Uint8     uint8       `json:"uint8"`
	Uint16    uint16      `json:"uint16"`
	Uint32    uint32      `json:"uint32"`
	Uint64    uint64      `json:"uint64"`
	Float32   float32     `json:"float32"`
	Float64   float64     `json:"float64"`

	// Slice of primitive types
	ArrayBoolean   []bool        `json:"array_boolean"`
	ArrayInterface []interface{} `json:"array_interface"`
	ArrayString    []string      `json:"array_string"`
	ArrayRandom    []any         `json:"array_random"`
	ArrayInt       []int         `json:"array_int"`
	ArrayInt8      []int8        `json:"array_int8"`
	ArrayInt16     []int16       `json:"array_int16"`
	ArrayInt32     []int32       `json:"array_int32"`
	ArrayInt64     []int64       `json:"array_int64"`
	ArrayUint      []uint        `json:"array_uint"`
	ArrayUint8     []uint8       `json:"array_uint8"`
	ArrayUint16    []uint16      `json:"array_uint16"`
	ArrayUint32    []uint32      `json:"array_uint32"`
	ArrayUint64    []uint64      `json:"array_uint64"`
	ArrayFloat32   []float32     `json:"array_float32"`
	ArrayFloat64   []float64     `json:"array_float64"`
}

type MapStruct struct {
	SSMap                    map[string]string                       `json:"ssmap"`
	ArrayOfMaps              []map[string]string                     `json:"array_of_maps"`
	MapOfMaps                map[string]map[string]string            `json:"map_of_maps"`
	MapStringInt             map[string]int                          `json:"map_string_int"`
	MapStringFloat           map[string]float64                      `json:"map_string_float"`
	MapStringInterface       map[string]interface{}                  `json:"map_string_interface"`
	MapStringAny             map[string]any                          `json:"map_string_any"`
	MapStringPrimitiveStruct map[string]PrimitiveStruct              `json:"map_string_primitive_struct"`
	MapOfMapsOfMaps          map[string]map[string]map[string]string `json:"map_of_maps_of_maps"`
	MapOfArrayOfMaps         map[string][]map[string]string          `json:"map_of_arrays_of_maps"`
}

type StructWithEmbeddedGenericStruct struct {
	SimpleStruct
	exstructpkg3.SingleGenericStructPkg3[exstructpkg2.SimpleStructPkg2, SimpleStruct]
	StringArray []string
	StructArray []string
}

type StructWithInlineStruct struct {
	InlineStructData struct {
		Test string `json:"test"`
	}
}
