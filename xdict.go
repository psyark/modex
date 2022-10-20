package modex

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"time"
)

// Dict はモデルが登録される辞書です
type Dict struct {
	items map[string]*xStruct
}

// Register はこの辞書にタイプを追加します
func (d *Dict) Register(refType reflect.Type) {
	if d.items == nil {
		d.items = map[string]*xStruct{}
	}

	// ポインタやスライスは型を辿る
	refType = getElementType(refType)

	// 特定のタイプは無視
	if refType == reflect.TypeOf(time.Time{}) {
		return
	}

	// 構造体でなければ無視
	if refType.Kind() != reflect.Struct {
		return
	}

	// 既に登録されていれば無視
	if _, ok := d.items[refType.Name()]; ok {
		return
	}

	// 登録
	d.items[refType.Name()] = &xStruct{refType}

	// 各フィールドのタイプも登録する
	for i := 0; i < refType.NumField(); i++ {
		sf := refType.Field(i)
		xf := xField{sf}
		if sf.Anonymous {
			d.Register(sf.Type)
		} else if xf.exportable() {
			if _, ok := sf.Tag.Lookup("modex"); !ok {
				d.Register(sf.Type)
			}
		}
	}
}

// Export は登録した全モデルをエクスポートします
func (d *Dict) Export(w io.Writer) {
	if d.items != nil {
		names := []string{}
		for name := range d.items {
			names = append(names, name)
		}
		sort.Strings(names)

		for _, name := range names {
			fmt.Fprintln(w, d.items[name])
		}
	}
}
