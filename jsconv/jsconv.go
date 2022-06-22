// Package jsconv provides conversions between Go and JS values.
package jsconv

import (
	"syscall/js"

	"github.com/blitz-frost/conv"
)

var global = js.Global()
var array = global.Get("Array")
var object = global.Get("Object")

var (
	To   func(*js.Value, any) error
	From func(any, js.Value) error
)

func arrayTo(dst *js.Value, src conv.Array) error {
	return arrayishTo(dst, src)
}

func arrayishTo(dst *js.Value, src conv.ArrayInterface) error {
	v := make([]any, src.Len())

	for i := range v {
		vi := js.Value{}
		if err := To(&vi, src.Index(i)); err != nil {
			return err
		}
		v[i] = vi
	}

	*dst = js.ValueOf(v)
	return nil
}

func directTo[t any](dst *js.Value, src t) error {
	*dst = js.ValueOf(src)
	return nil
}

func mapTo(dst *js.Value, src conv.Map) error {
	v := make(map[string]any, src.Len())

	iter := src.Range()
	for iter.Next() {
		k := iter.Key().(string)
		vk := js.Value{}
		if err := To(&vk, iter.Value()); err != nil {
			return err
		}
		v[k] = vk
	}

	*dst = js.ValueOf(v)
	return nil
}

func nilTo(dst *js.Value, src conv.Nil) error {
	*dst = js.Null()
	return nil
}

func sliceTo(dst *js.Value, src conv.Slice) error {
	return arrayishTo(dst, src)
}

func init() {
	scheme := conv.MakeScheme(js.Value{})
	scheme.Load(arrayTo)
	scheme.Load(directTo[bool])
	scheme.Load(directTo[float64])
	scheme.Load(directTo[int64])
	scheme.Load(directTo[string])
	scheme.Load(directTo[uint64])
	scheme.Load(mapTo)
	scheme.Load(nilTo)
	scheme.Load(sliceTo)
	scheme.Build(&To)
}
