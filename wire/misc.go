// This file defines conversions for simple types.

package wire

import (
	"reflect"
	"unsafe"
)

var simpleKinds map[reflect.Kind]bool = map[reflect.Kind]bool{
	reflect.Bool:       true,
	reflect.Int:        true,
	reflect.Int8:       true,
	reflect.Int16:      true,
	reflect.Int32:      true,
	reflect.Int64:      true,
	reflect.Uint:       true,
	reflect.Uint8:      true,
	reflect.Uint16:     true,
	reflect.Uint32:     true,
	reflect.Uint64:     true,
	reflect.Float32:    true,
	reflect.Float64:    true,
	reflect.Complex64:  true,
	reflect.Complex128: true,
}

func marshalSimple(ptr uintptr, size int) []byte {
	h := reflect.SliceHeader{
		Data: ptr,
		Cap:  size,
		Len:  size,
	}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func marshalInt(i int) []byte {
	return (*(*[uintSize]byte)(unsafe.Pointer(&i)))[:]
}

// addressable returns back v if it is addressable, as per the reflect package, or a copy that is.
func addressable(v reflect.Value) reflect.Value {
	if v.CanAddr() {
		return v
	}

	r := reflect.New(v.Type()).Elem()
	r.Set(v)
	return r
}
