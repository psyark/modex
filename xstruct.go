package modex

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

// xStruct はDictのエントリーで、KindがStructであるTypeを持ちます
type xStruct struct {
	reflect.Type
}

// String はこのxStructの文字列表現をTypeScript interfaceで返します
func (s *xStruct) String() string {
	buf := bytes.NewBufferString("")
	fmt.Fprintf(buf, "%v {\n", s.getSignature())
	for i := 0; i < s.Type.NumField(); i++ {
		f := &xField{s.Type.Field(i)}
		f.export(buf)
	}
	fmt.Fprint(buf, "}\n")
	return buf.String()
}

func (s *xStruct) getSignature() string {
	extendNames := []string{}
	for i := 0; i < s.Type.NumField(); i++ {
		refFld := s.Type.Field(i)
		if refFld.Anonymous {
			switch refFld.Type.Kind() {
			case reflect.Struct:
				extendNames = append(extendNames, refFld.Type.Name())
			case reflect.Ptr:
				extendNames = append(extendNames, refFld.Type.Elem().Name())
			}
		}
	}

	signature := fmt.Sprintf("export interface %v", s.Type.Name())
	if len(extendNames) != 0 {
		signature += fmt.Sprintf(" extends %v", strings.Join(extendNames, ", "))
	}
	return signature
}
