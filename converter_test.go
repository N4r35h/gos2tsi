package gos2tsi

import (
	"testing"

	"github.com/N4r35h/gos2tsi/examplestructs"
)

var c *Converter = New()

func TestSimpleStruct1(t *testing.T) {
	ps := c.ParseStruct(examplestructs.SimpleStruct1{})
	op := c.GetStructAsInterfaceString(ps)
	expected := `export interface SimpleStruct1 {
string_field: string
boolean_field: boolean
uint_field: number
int_field: number
float64_field: number
FieldWOJSONTag: string
}`
	if op != expected {
		t.Errorf(expected)
		t.Errorf(op)
	}
}

func TestEmbededStruct(t *testing.T) {
	ps := c.ParseStruct(examplestructs.StructWithEmbeding{})
	op := c.GetStructAsInterfaceString(ps)
	expected := `export interface StructWithEmbeding {
id: number
level: string
}`
	if op != expected {
		t.Errorf(expected)
		t.Errorf(op)
	}
}

func TestStructWithFieldStruct(t *testing.T) {
	ps := c.ParseStruct(examplestructs.StructWithFieldStruct{})
	op := c.GetStructAsInterfaceString(ps)
	expected := `export interface StructWithFieldStruct {
struct_field: SimpleStruct
}`
	if op != expected {
		t.Errorf(expected)
		t.Errorf(op)
	}
}

func TestStructWithFieldSlice(t *testing.T) {
	ps := c.ParseStruct(examplestructs.StructWithFieldSlice{})
	op := c.GetStructAsInterfaceString(ps)
	expected := `export interface StructWithFieldSlice {
slice_field: string[]
}`
	if op != expected {
		t.Errorf(expected)
		t.Errorf(op)
	}
}

func TestStructWithFieldStructSlice(t *testing.T) {
	ps := c.ParseStruct(examplestructs.StructWithFieldStructSlice{})
	op := c.GetStructAsInterfaceString(ps)
	expected := `export interface StructWithFieldStructSlice {
struct_slice_field: SimpleStruct[]
}`
	if op != expected {
		t.Errorf(expected)
		t.Errorf(op)
	}
}

func TestSliceOfSimpleStruct(t *testing.T) {
	ps := c.ParseStruct([]examplestructs.SimpleStruct{})
	op := c.GetStructAsInterfaceString(ps)
	expected := `export interface SimpleStruct {
test: string
}`
	if op != expected {
		t.Errorf(expected)
		t.Errorf(op)
	}
	if ps.IsSlice != 1 {
		t.Errorf("ps.IsSlice must be true")
	}
}

func TestSingleGenericStruct1(t *testing.T) {
	ps := c.ParseStruct(examplestructs.SingleGenericStruct[examplestructs.SimpleStruct]{})
	op := c.GetStructAsInterfaceString(ps)
	expected := `export interface SingleGenericStruct<T> {
data: T[]
pagination_info: SimpleStruct
}`
	if op != expected {
		t.Errorf(expected)
		t.Errorf(op)
	}
}

func TestSingleGenericStruct2(t *testing.T) {
	ps := c.ParseStruct(examplestructs.SingleGenericStruct[examplestructs.SimpleStruct1]{})
	op := c.GetStructAsInterfaceString(ps)
	expected := `export interface SingleGenericStruct<T> {
data: T[]
pagination_info: SimpleStruct
}`
	if op != expected {
		t.Errorf(expected)
		t.Errorf(op)
	}
}

func TestMultiGenericStruct(t *testing.T) {
	ps := c.ParseStruct(examplestructs.MultiGenericStruct[examplestructs.SimpleStruct, examplestructs.SimpleStruct1]{})
	op := c.GetStructAsInterfaceString(ps)
	expected := `export interface MultiGenericStruct<T, U> {
data: T[]
data2: U[]
pagination_info: SimpleStruct
}`
	if op != expected {
		t.Errorf(expected)
		t.Errorf(op)
	}
}

func TestStructWithCommentOnTop(t *testing.T) {
	ps := c.ParseStruct(examplestructs.StructWithCommentOnTop{})
	op := c.GetStructAsInterfaceString(ps)
	expected := `
/**
Test
Test 123
 Test 123
*/
export interface StructWithCommentOnTop {
test: string
}`
	if op != expected {
		t.Errorf(expected)
		t.Errorf(op)
	}
}

func TestStructWithCustomTSType(t *testing.T) {
	ps := c.ParseStruct(examplestructs.StructWithCustomTSType{})
	op := c.GetStructAsInterfaceString(ps)
	expected := `export interface StructWithCustomTSType {
custom_ts_type: string
}`
	if op != expected {
		t.Errorf(expected)
		t.Errorf(op)
	}
}

func TestPrimitiveStruct(t *testing.T) {
	ps := c.ParseStruct(examplestructs.PrimitiveStruct{})
	t.Logf("%+v \n", ps)
	op := c.GetStructAsInterfaceString(ps)
	expected := `export interface PrimitiveStruct {
boolean: boolean
interface: any
string: string
random: any
int: number
int8: number
int16: number
int32: number
int64: number
uint: number
uint8: number
uint16: number
uint32: number
uint64: number
float32: number
float64: number
array_boolean: boolean[]
array_interface: any[]
array_string: string[]
array_random: any[]
array_int: number[]
array_int8: number[]
array_int16: number[]
array_int32: number[]
array_int64: number[]
array_uint: number[]
array_uint8: number[]
array_uint16: number[]
array_uint32: number[]
array_uint64: number[]
array_float32: number[]
array_float64: number[]
}`
	if op != expected {
		t.Errorf(expected)
		t.Errorf(op)
	}
}

func TestStructWithMaps(t *testing.T) {

	ps := c.ParseStruct(examplestructs.MapStruct{})
	op := c.GetStructAsInterfaceString(ps)
	expected := `export interface MapStruct {
ssmap: {[key: string]: string}
array_of_maps: {[key: string]: string}[]
map_of_maps: {[key: string]: {[key: string]: string}}
map_string_int: {[key: string]: number}
map_string_float: {[key: string]: number}
map_string_interface: {[key: string]: any}
map_string_any: {[key: string]: any}
map_string_primitive_struct: PrimitiveStruct
map_of_maps_of_maps: {[key: string]: {[key: string]: {[key: string]: string}}}
map_of_arrays_of_maps: {[key: string]: {[key: string]: string}[]}
}`
	if op != expected {
		t.Errorf(expected)
		t.Errorf(op)
	}
}

// var allStructs []any = []any{
// 	examplestructs.EmbedableStruct{},
// 	examplestructs.StructWithEmbeding{},
// 	examplestructs.MapStruct{},
// 	examplestructs.MultiGenericStruct[any, any]{},
// 	examplestructs.PrimitiveStruct{},
// 	examplestructs.SimpleStruct{},
// 	examplestructs.SimpleStruct1{},
// 	examplestructs.StructWithFieldStruct{},
// 	examplestructs.StructWithFieldSlice{},
// 	examplestructs.StructWithFieldStructSlice{},
// 	examplestructs.SingleGenericStruct[any]{},
// 	examplestructs.StructWithCommentOnTop{},
// 	examplestructs.StructWithCustomTSType{},
// }

// func BenchmarkGOSTSI(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		for _, v := range allStructs {
// 			ps := c.ParseStruct(v)
// 			c.GetStructAsInterfaceString(ps)
// 		}
// 	}
// }

// func BenchmarkGUT(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		for _, v := range allStructs {
// 			gut.Convert(v)
// 		}
// 	}
// }

func TestStructWithEmddedGenericStructs(t *testing.T) {

	ps := c.ParseStruct(examplestructs.StructWithEmbeddedGenericStruct{})
	op := c.GetStructAsInterfaceString(ps)
	expected := `export interface StructWithEmbeddedGenericStruct {
test: string
data: SimpleStructPkg2[]
data2: {[key: string]: SimpleStruct}
pagination_info: string
StringArray: string[]
StructArray: string[]
}`
	if op != expected {
		t.Errorf(expected)
		t.Errorf(op)
	}
}

func TestStructWithGenericPrimitiveSlice(t *testing.T) {
	ps := c.ParseStruct(examplestructs.SingleGenericStruct[[]string]{})
	op := c.GetStructAsInterfaceString(ps)
	expected := `export interface SingleGenericStruct<T> {
data: T[]
pagination_info: SimpleStruct
}`
	if op != expected {
		t.Errorf(expected)
		t.Errorf(op)
	}
}

func TestStructWithGenericStructSlice(t *testing.T) {
	ps := c.ParseStruct(examplestructs.SingleGenericStruct[[]examplestructs.SimpleStruct]{})
	op := c.GetStructAsInterfaceString(ps)
	expected := `export interface SingleGenericStruct<T> {
data: T[]
pagination_info: SimpleStruct
}`
	if op != expected {
		t.Errorf(expected)
		t.Errorf(op)
	}
}

func TestStructWithInlineStruct(t *testing.T) {

	ps := c.ParseStruct(examplestructs.StructWithInlineStruct{})
	op := c.GetStructAsInterfaceString(ps)
	expected := `export interface StructWithInlineStruct {
InlineStructData: struct{Test string "json:\"test\""}
}`
	if op != expected {
		t.Errorf(expected)
		t.Errorf(op)
	}
}

func Test3DSlice(t *testing.T) {
	ps := c.ParseStruct(examplestructs.Slice3D{})
	op := c.GetStructAsInterfaceString(ps)
	expected := `export interface Slice3D {
String3DSlice: string[][]
Struct3DSlice: SimpleStruct[][]
}`
	if op != expected {
		t.Errorf(expected)
		t.Errorf(op)
	}
}

// if a given field has optional tag or omitempty tag, it should be optional in the ts interface
func TestOptionalField(t *testing.T) {
	ps := c.ParseStruct(examplestructs.StructWithOptionalField{})
	op := c.GetStructAsInterfaceString(ps)
	expected := `export interface StructWithOptionalField {
required_field: string
omitempty_field?: string
optional_field?: string
}`

	if op != expected {
		t.Errorf(expected)
		t.Errorf(op)
	}
}

func TestPointerFields(t *testing.T) {
	ps := c.ParseStruct(examplestructs.StructWithPointers{})
	op := c.GetStructAsInterfaceString(ps)
	expected := `export interface StructWithPointers {
int_field: number
bool_field: boolean
array_field: []string
entity_x: StructWithOptionalField
}`
	if op != expected {
		t.Errorf(expected)
		t.Errorf(op)
	}
}
