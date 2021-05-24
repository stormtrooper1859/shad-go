// +build !solution

package illegal

import "unsafe"

func StringFromBytes(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
