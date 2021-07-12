// +build !solution

package illegal

import (
	"reflect"
	"unsafe"
)

func SetPrivateField(obj interface{}, name string, value interface{}) {
	reflObj := reflect.ValueOf(obj)
	reflValue := reflect.ValueOf(value)
	reflField, _ := reflObj.Type().Elem().FieldByName(name)
	fieldPointer := (unsafe.Pointer)(reflObj.Pointer() + reflField.Offset)

	r := reflect.NewAt(reflValue.Type(), fieldPointer)
	r.Elem().Set(reflValue)
}
