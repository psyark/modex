package modex

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// タイプがElemを持つ限り、そのElemのタイプを辿って返します
func getElementType(t reflect.Type) (et reflect.Type) {
	defer func() {
		recover()
	}()

	et = t
	for {
		et = et.Elem()
	}
}

func toTsType(refTyp reflect.Type) string {
	switch refTyp {
	case reflect.TypeOf(time.Time{}):
		return "string"
	}

	switch refTyp.Kind() {
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return "number"
	case reflect.String:
		return "string"
	case reflect.Struct:
		return refTyp.Name()
	case reflect.Array, reflect.Slice:
		elemType := toTsType(refTyp.Elem())
		if strings.HasSuffix(elemType, " | null") {
			elemType = "(" + elemType + ")"
		}
		return elemType + "[]"
	case reflect.Ptr:
		return toTsType(refTyp.Elem()) + " | null"
	case reflect.Map:
		return fmt.Sprintf("{[key: %v]: %v}", toTsType(refTyp.Key()), toTsType(refTyp.Elem()))
	case reflect.Interface:
		return "any"
	// case reflect.Uintptr:
	// case reflect.Complex64:
	// case reflect.Complex128:
	// case reflect.Chan:
	// case reflect.Func:
	// case reflect.UnsafePointer:
	default:
		panic(fmt.Sprintf("type=%v, kind=%v", refTyp, refTyp.Kind()))
	}
}
