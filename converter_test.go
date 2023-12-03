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
	if ps.IsSlice != true {
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
