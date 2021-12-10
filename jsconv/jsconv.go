// Package jsconv provides conversion between Go and JS values.
package jsconv

import (
	"errors"
	"reflect"
	"strings"
	"syscall/js"
)

var global = js.Global()
var array = global.Get("Array")
var object = global.Get("Object")

var stringType = reflect.TypeOf("")

// these types can be directly converted with js.ValueOf
var (
	directSlice     = reflect.TypeOf([]interface{}(nil))
	directMap       = reflect.TypeOf(map[string]interface{}(nil))
	directInterface = reflect.TypeOf((*js.Wrapper)(nil)).Elem()

	directTypes = map[reflect.Type]struct{}{
		reflect.TypeOf(false):      struct{}{},
		reflect.TypeOf(int(0)):     struct{}{},
		reflect.TypeOf(int8(0)):    struct{}{},
		reflect.TypeOf(int16(0)):   struct{}{},
		reflect.TypeOf(int32(0)):   struct{}{},
		reflect.TypeOf(int64(0)):   struct{}{},
		reflect.TypeOf(uint(0)):    struct{}{},
		reflect.TypeOf(uint8(0)):   struct{}{},
		reflect.TypeOf(uint16(0)):  struct{}{},
		reflect.TypeOf(uint32(0)):  struct{}{},
		reflect.TypeOf(uint64(0)):  struct{}{},
		reflect.TypeOf(uintptr(0)): struct{}{},
		reflect.TypeOf(float32(0)): struct{}{},
		reflect.TypeOf(float64(0)): struct{}{},
		reflect.TypeOf(""):         struct{}{},
		directSlice:                struct{}{},
		directMap:                  struct{}{},
	}
)

// checking simple kinds with a map is faster than trying CanConvert on all direct types
var indirectKinds = make(map[reflect.Kind]reflect.Type)

func init() {
	for typ, _ := range directTypes {
		indirectKinds[typ.Kind()] = typ
	}

	// exclude []interface{} and map[string]interface{}
	delete(indirectKinds, reflect.Slice)
	delete(indirectKinds, reflect.Map)
}

// eval returns the kind and direct convertible (using js.ValueOf) type of v (or nil if it doesn't correspond to one).
func eval(v reflect.Value) (reflect.Kind, reflect.Type) {
	k := v.Kind()
	if directType, ok := indirectKinds[k]; ok {
		return k, directType
	}
	if v.CanConvert(directSlice) {
		return k, directSlice
	}
	if v.CanConvert(directMap) {
		return k, directMap
	}
	if v.CanConvert(directInterface) {
		return k, directInterface
	}
	return k, nil
}

var integerKinds = map[reflect.Kind]struct{}{
	reflect.Int:    struct{}{},
	reflect.Int8:   struct{}{},
	reflect.Int16:  struct{}{},
	reflect.Int32:  struct{}{},
	reflect.Int64:  struct{}{},
	reflect.Uint:   struct{}{},
	reflect.Uint8:  struct{}{},
	reflect.Uint16: struct{}{},
	reflect.Uint32: struct{}{},
	reflect.Uint64: struct{}{},
}

// JSValue can be used to defer JS to Go conversion using the "From" methods.
type JSValue js.Value

func (x JSValue) JSValue() js.Value {
	return js.Value(x)
}

var jsValueType = reflect.TypeOf(JSValue{})

// From converts the source js value, storing it into the destination go pointer.
// dst must be a pointer to an appropriate type.
func From(dst interface{}, src js.Wrapper) error {
	if dst == nil {
		return errors.New("nil destination")
	}

	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr {
		return errors.New("destination not a pointer")
	}

	return FromValue(v.Elem(), src)
}

// FromValue is the underling implementation of From.
// dst must be a settable destination value or a pointer to an appropriate value.
func FromValue(dst reflect.Value, src js.Wrapper) error {
	srcVal := src.JSValue()

	v := reflect.Indirect(dst)
	t := v.Type()
	k := v.Kind()

	// check for defered conversion
	if t.ConvertibleTo(jsValueType) {
		v.Set(reflect.ValueOf(srcVal).Convert(t))
		return nil
	}

	switch srcVal.Type() {
	case js.TypeUndefined:
		fallthrough
	case js.TypeNull:
		v.Set(reflect.Zero(t))
	case js.TypeBoolean:
		if k != reflect.Bool {
			return errors.New("destination not bool pointer")
		}
		v.SetBool(srcVal.Bool())
	case js.TypeNumber:
		if _, ok := integerKinds[k]; ok {
			v.SetInt(int64(srcVal.Int()))
		} else if k == reflect.Float32 || k == reflect.Float64 {
			v.SetFloat(srcVal.Float())
		} else {
			return errors.New("destination not number pointer")
		}
	case js.TypeString:
		if k != reflect.String {
			return errors.New("destination not string pointer")
		}
		v.SetString(srcVal.String())
	case js.TypeObject:
		if array.Call("isArray", src).Bool() {
			n := srcVal.Length()
			switch k {
			case reflect.Array:
				if v.Len() < n {
					return errors.New("destination is pointer to array of insufficient length")
				}
				v.Set(reflect.Zero(t)) // clear array
			case reflect.Slice:
				if v.Cap() < n {
					v.Set(reflect.MakeSlice(t, n, n))
				} else if v.Len() < n {
					v.SetLen(n)
				}
			default:
				return errors.New("destination not array or slice")
			}

			for i := 0; i < n; i++ {
				if err := FromValue(v.Index(i), srcVal.Index(i)); err != nil {
					return err
				}
			}
		}

		switch k {
		case reflect.Struct:
			fields := reflect.VisibleFields(t)
			for _, field := range fields {
				s := unexported(field.Name)
				if err := FromValue(v.FieldByIndex(field.Index), srcVal.Get(s)); err != nil {
					return err
				}
			}
		case reflect.Map:
			keyType := t.Key()
			if keyType.Kind() != reflect.String {
				return errors.New("destination not map with string keys")
			}
			valType := t.Elem()

			v.Set(reflect.MakeMap(t)) // clear map

			keys := object.Call("keys", src)
			n := keys.Length()
			for i := 0; i < n; i++ {
				keyStr := keys.Index(i).String() // key string

				key := reflect.New(keyType).Elem() // key value
				key.SetString(keyStr)

				val := reflect.New(valType).Elem()
				if err := FromValue(val, srcVal.Get(keyStr)); err != nil {
					return err
				}

				v.SetMapIndex(key, val)
			}
		default:
			return errors.New("destination not struct or map")
		}
	default:
		return errors.New("invalid js value")
	}

	return nil

}

// toDirect returns a value that is directly convertible using js.ValueOf
func toDirect(v reflect.Value) (interface{}, error) {
	v = reflect.Indirect(v)
	t := reflect.TypeOf(v)

	if _, ok := directTypes[t]; ok {
		return v.Interface(), nil
	}

	k, direct := eval(v)
	if direct != nil {
		return v.Convert(direct).Interface(), nil
	}

	var o interface{}
	switch k {
	case reflect.Slice:
		if v.IsNil() {
			return nil, nil
		}

		fallthrough
	case reflect.Array:
		oo := make([]interface{}, v.Len())
		for i := range oo {
			elem, err := toDirect(v.Index(i))
			if err != nil {
				return nil, err
			}
			oo[i] = elem
		}
		o = oo
	case reflect.Map:
		if v.IsNil() {
			return nil, nil
		}

		if t.Key().Kind() != reflect.String {
			return nil, errors.New("map key not string")
		}

		oo := make(map[string]interface{})
		iter := v.MapRange()
		for iter.Next() {
			val, err := toDirect(iter.Value())
			if err != nil {
				return nil, err
			}
			key := iter.Key().Convert(stringType).String()
			oo[key] = val
		}
		o = oo
	case reflect.Struct: // convert exported struct fields to map
		oo := make(map[string]interface{})
		fields := reflect.VisibleFields(t)
		for _, field := range fields {
			val, err := toDirect(v.FieldByIndex(field.Index))
			if err != nil {
				return nil, err
			}
			oo[unexported(field.Name)] = val
		}
		o = oo
	default:
		return nil, errors.New("invalid value")
	}

	return o, nil
}

// To converts the input Go value to JS.
// Pointers are flattened.
// Only maps with string keys are allowed.
func To(src interface{}) (js.Value, error) {
	// reflect.ValueOf(nil) produces an invalid Value
	if src == nil {
		return js.Null(), nil
	}

	return ToValue(reflect.ValueOf(src))
}

// ToValue is the reflect version of "To".
func ToValue(v reflect.Value) (js.Value, error) {
	direct, err := toDirect(v)
	if err != nil {
		return js.Undefined(), err
	}
	return js.ValueOf(direct), nil
}

// unexported returns the unexported version of the string s (lower case first letter)
func unexported(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}
