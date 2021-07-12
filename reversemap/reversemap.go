// +build !solution

package reversemap

import (
	"reflect"
)

func ReverseMap(forward interface{}) interface{} {
	v := reflect.ValueOf(forward)
	if v.Kind() != reflect.Map {
		panic("not a map")
	}

	keyType := v.Type().Key()
	valueType := v.Type().Elem()

	result := reflect.MakeMap(reflect.MapOf(valueType, keyType))
	for _, key := range v.MapKeys() {
		value := v.MapIndex(key)
		result.SetMapIndex(value, key)
	}

	return result.Interface()
}
