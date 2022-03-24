//go:build !appengine && !nounsafe

package h1

import (
	"reflect"
	"unsafe"
)

//nolint
func stringToBytes(s string) []byte {
	//#nosec
	return unsafe.Slice((*byte)(unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&s)).Data)), len(s))
}

func bytesToString(b []byte) string {
	//#nosec
	return *(*string)(unsafe.Pointer(&b))
}
