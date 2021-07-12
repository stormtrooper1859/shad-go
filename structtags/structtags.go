// +build !solution

package structtags

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

var cache sync.Map

func Unpack(req *http.Request, ptr interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	var fields map[string]int

	v := reflect.ValueOf(ptr).Elem()
	fieldsAny, contain := cache.Load(v.Type())

	if !contain {
		fields = make(map[string]int)

		for i := 0; i < v.NumField(); i++ {
			fieldInfo := v.Type().Field(i)
			tag := fieldInfo.Tag
			name := tag.Get("http")
			if name == "" {
				name = strings.ToLower(fieldInfo.Name)
			}
			fields[name] = i
		}

		cache.Store(v.Type(), fields)
	} else {
		fields = fieldsAny.(map[string]int)
	}

	for name, values := range req.Form {
		fieldNumber, ok := fields[name]
		if !ok {
			continue
		}

		f := v.Field(fieldNumber)

		for _, value := range values {
			if f.Kind() == reflect.Slice {
				elem := reflect.New(f.Type().Elem()).Elem()
				if err := populate(elem, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
				f.Set(reflect.Append(f, elem))
			} else {
				if err := populate(f, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
			}
		}
	}
	return nil
}

func populate(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)

	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)

	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(b)

	default:
		return fmt.Errorf("unsupported kind %s", v.Type())
	}
	return nil
}
