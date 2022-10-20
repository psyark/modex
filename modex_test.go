package modex

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
	"time"
	"unsafe"
)

func TestModex(t *testing.T) {
	testCases := []struct {
		name string
		args interface{}
		want string
	}{
		{"primitives", &Primitives{}, primitivesResult},
		{"slice", &Slice{}, sliceResult},
		{"complex", &Complex{}, complexResult},
		{"jsontag", &JSONTag{}, jsontagResult},
		{"modextag", &ModexTag{}, modextagResult},
		{"unexport", &Unexport{}, unexportResult},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dict := Dict{}
			dict.Register(reflect.TypeOf(tc.args))

			buf := bytes.NewBufferString("")
			dict.Export(buf)
			got := buf.String()

			if strings.TrimSpace(tc.want) != strings.TrimSpace(got) {
				t.Errorf("got:\n%v, want:\n%v", got, tc.want)
			}
		})
	}
}

type Primitives struct {
	Bool    bool
	Int     int
	Int8    int8
	Int16   int16
	Int32   int32
	Int64   int64
	Uint    uint
	Uint8   uint8
	Uint16  uint16
	Uint32  uint32
	Uint64  uint64
	Float32 float32
	Float64 float64
	String  string
}

const primitivesResult = `
export interface Primitives {
    Bool: boolean;
    Int: number;
    Int8: number;
    Int16: number;
    Int32: number;
    Int64: number;
    Uint: number;
    Uint8: number;
    Uint16: number;
    Uint32: number;
    Uint64: number;
    Float32: number;
    Float64: number;
    String: string;
}
`

type JSONTag struct {
	WithName    string `json:"hoge"`
	NoExport    string `json:"-"`
	OmitEmpty   string `json:"fuga,omitempty"`
	OmitEmpty2  string `json:",omitempty"`
	InvalidName string `json:".,omitempty"`
}

const jsontagResult = `
export interface JSONTag {
    hoge: string;
    fuga?: string;
    OmitEmpty2?: string;
    '.'?: string;
}
`

type ModexTag struct {
	Time   time.Time
	Time2  time.Time          `modex:"number"`
	String ModexTagUnexported `modex:"string"`
}
type ModexTagUnexported struct{}

const modextagResult = `
export interface ModexTag {
    Time: string;
    Time2: number;
    String: string;
}
`

type Slice struct {
	Array        [2]string
	Slice        []string
	PtrSlice     []*string
	SlicePtr     *[]string
	NestedSlice  [][]string
	NestedSlice2 [][][]string
}

const sliceResult = `
export interface Slice {
    Array: string[];
    Slice: string[];
    PtrSlice: (string | null)[];
    SlicePtr: string[] | null;
    NestedSlice: string[][];
    NestedSlice2: string[][][];
}
`

type ComplexBaseBase struct{}
type ComplexBase struct {
	*ComplexBaseBase
	Base string
}
type Complex struct {
	ComplexBase
	Slice []ComplexBase
	Ptr   *ComplexBase
	Map   map[string]ComplexBase
}

const complexResult = `
export interface Complex extends ComplexBase {
    Slice: ComplexBase[];
    Ptr: ComplexBase | null;
    Map: {[key: string]: ComplexBase};
}

export interface ComplexBase extends ComplexBaseBase {
    Base: string;
}

export interface ComplexBaseBase {
}
`

type Unexport struct {
	Time      time.Time
	Interface interface{}
	Func      func() string
	Chan      chan int
	// ChanPtr   *(chan int)
	Unsafe  unsafe.Pointer
	private string
}

const unexportResult = `
export interface Unexport {
    Time: string;
    Interface: any;
}
`
