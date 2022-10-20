package modex

import (
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strings"
	"unicode"
)

var validNameRegex *regexp.Regexp

func init() {
	validNameRegex = regexp.MustCompile(`^([_\w]+)$`)
}

// xField はStructFieldのラッパーです
type xField struct {
	reflect.StructField
}

func (f *xField) exportable() bool {
	field := f.StructField

	// 埋め込みフィールドは（interfaceのメンバーとしては）エクスポートしない
	if field.Anonymous {
		return false
	}

	// 特定の型はエクスポートしない
	// TODO ElemがChanに対応
	if field.Type.Kind() == reflect.Chan {
		return false
	}

	switch getElementType(field.Type).Kind() {
	case reflect.Func:
		return false
	case reflect.UnsafePointer:
		return false
	}

	// 非公開メンバーはエクスポートしない
	if !unicode.IsUpper([]rune(field.Name)[0]) {
		return false
	}

	name, _ := f.parseJSONTag()
	return name != "-"
}

// parseJSONTag はJSONタグをパースし、指定された別名（または"-"）およびomitemptyの有無を返します
func (f *xField) parseJSONTag() (string, bool) {
	if jsonTag, ok := f.StructField.Tag.Lookup("json"); ok {
		jsonTagItems := strings.SplitN(jsonTag, ",", 2)
		return jsonTagItems[0], strings.Contains(jsonTag+",", ",omitempty,")
	}
	return "", false
}

// getTsName はJSONタグをパースし、TypeScriptコードに出力できる（必要なら引用符付きの）名前を返します
func (f *xField) getTsName() string {
	name, omitempty := f.parseJSONTag()
	if name == "" {
		name = f.StructField.Name
	} else if !validNameRegex.MatchString(name) {
		name = fmt.Sprintf("'%v'", name)
	}
	if omitempty {
		name += "?"
	}
	return name
}

func (f *xField) export(w io.Writer) {
	sf := f.StructField
	if f.exportable() {
		tsType, ok := sf.Tag.Lookup("modex")
		if !ok {
			tsType = toTsType(sf.Type)
		}
		fmt.Fprintf(w, "    %v: %v;\n", f.getTsName(), tsType)
	}
}
