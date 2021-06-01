/*
Package wire provides Go native binary serialization.

Decoding is based on providing an initialized pointer, as in the standard library encoding packages.
However, one of the principles encoding in this package is for data to be self describing.
This would potentially allow extensions that can reconstruct data without having to use an existing pointer, effectively providing runtime reflect capabilities across programs.

Encoding details

All encoded values follow the format:
	kindID value
Basic numeric types (bool, uint, int, complex) directly use their in-memory representation.

Arrays and slices use different IDs, but their values are encoded identically: one int representing the length, followed by that many recursively encoded values.

Maps start with one int for the number of key-value pairs, followed by each key and value individually encoded.

Structs must only have exported fields. Field names are not transmitted or checked, the receiving side must use a struct type with the right field number, type and order.
The encoded form is: one int representing the number of fields, followed by each field encoded individually.

No other types are currently supported. Notably, pointers and interfaces are not supported. All used types must be concrete and flat (not be or contain references).
*/
package wire

import (
	"errors"
	"io"
	"reflect"
	"unsafe"
)

type T struct {
	a int
	b int
}

const (
	uintSize = 4 << (^uint(0) >> 32 & 1) // 4 or 8
)

// A size encodes counts or byte lengths.
// It is an int type for compatibility with the built-in Go len function, but is not used for negative values.
// It is meant to strike a balance between processing speed and encoding size.
// It can hold positive numbers up to 2^56-1. Encodes as follows:
//
// First byte - number of additional bytes it occupies (up to 7)
//
// Rest - actual number as unsigned integer
//
// currently unused
type size int

func (x size) Encode(w io.Writer) error {
	x <<= 8
	b := *(*[uintSize]byte)(unsafe.Pointer(&x))
	i := byte(2)
	if b[7] != 0 {
		i = 8
	} else {
		for b[i] != 0 {
			i++
		}
	}
	b[0] = i - 1
	_, err := w.Write(b[:i])
	return err
}

func (x *size) Decode(r io.Reader) error {
	var b [uintSize + 1]byte
	if _, err := r.Read(b[:1]); err != nil {
		return err
	}
	if b[0] > 7 {
		return errors.New("size too large")
	}
	if _, err := r.Read(b[1 : b[0]+1]); err != nil {
		return err
	}
	*x = *(*size)(unsafe.Pointer(&b))
	*x = size(uint(*x) >> 8) // if the left shift operand is an integer, the shift is arithmetic (/2) instead of logical

	return nil
}

type Encoder struct {
	dst io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

// encodeInt is a convenience function to encode int values
func (x *Encoder) encodeInt(i int) error {
	_, err := x.dst.Write(marshalInt(i))
	return err
}

func (x *Encoder) Encode(v interface{}) error {
	return x.EncodeValue(reflect.ValueOf(v))
}

func (x *Encoder) EncodeValue(v reflect.Value) error {
	k := v.Kind()
	t := v.Type()

	if _, err := x.dst.Write([]byte{byte(k)}); err != nil {
		return err
	}
	switch {
	case simpleKinds[k]:
		v = addressable(v)
		_, err := x.dst.Write(marshalSimple(v.UnsafeAddr(), int(t.Size())))
		return err

	case k == reflect.Array || k == reflect.Slice:
		v = addressable(v)

		n := v.Len()

		if err := x.encodeInt(n); err != nil {
			return err
		}
		for i := 0; i < n; i++ {
			if err := x.EncodeValue(v.Index(i)); err != nil {
				return err
			}
		}

	case k == reflect.String:
		// write size
		if err := x.encodeInt(v.Len()); err != nil {
			return err
		}

		// convert string to byte slice and write it
		if _, err := x.dst.Write([]byte(v.Interface().(string))); err != nil {
			return err
		}

	case k == reflect.Map:
		n := v.Len()

		// write size
		if err := x.encodeInt(n); err != nil {
			return err
		}

		iter := v.MapRange()

		// encode each key value pair
		for i := 0; i < n; i++ {
			iter.Next()
			if err := x.EncodeValue(iter.Key()); err != nil {
				return err
			}
			if err := x.EncodeValue(iter.Value()); err != nil {
				return err
			}
		}

	case k == reflect.Struct:
		// currently not transmitting field names

		// write field count
		n := v.NumField()
		if err := x.encodeInt(n); err != nil {
			return err
		}

		// encode each field
		for i := 0; i < n; i++ {
			if err := x.EncodeValue(v.Field(i)); err != nil {
				return err
			}
		}

	default:
		return errors.New("unsupported type")
	}

	return nil
}

type Decoder struct {
	src io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r}
}

// decodeInt reads the next encoded value as an int
func (x *Decoder) decodeInt() (int, error) {
	r := make([]byte, uintSize)
	_, err := x.src.Read(r)
	return *(*int)(unsafe.Pointer(&r[0])), err
}

func (x *Decoder) Decode(v interface{}) error {
	return x.DecodeValue(reflect.ValueOf(v))
}

func (x *Decoder) DecodeValue(v reflect.Value) error {
	if v.Kind() != reflect.Ptr {
		return errors.New("not a pointer")
	}
	if v.IsNil() {
		return errors.New("nil pointer")
	}
	v = v.Elem()

	r := make([]byte, 1)
	if _, err := x.src.Read(r); err != nil {
		return err
	}
	k := reflect.Kind(r[0])
	if k != v.Kind() {
		return errors.New("incompatible value")
	}

	t := v.Type()
	switch {
	case simpleKinds[k]:
		r = make([]byte, int(t.Size()))
		if _, err := x.src.Read(r); err != nil {
			return err
		}

		p := reflect.NewAt(t, unsafe.Pointer(&r[0]))
		v.Set(p.Elem())

	case k == reflect.Array:
		// check length
		n, err := x.decodeInt()
		if err != nil {
			return err
		}
		if n != t.Len() {
			return errors.New("missmatching array length")
		}

		// populate array
		for i := 0; i < n; i++ {
			if err := x.DecodeValue(v.Index(i).Addr()); err != nil {
				return err
			}
		}

	case k == reflect.Slice:
		// get length
		n, err := x.decodeInt()
		if err != nil {
			return err
		}

		// make slice and populate it
		v.Set(reflect.MakeSlice(t, n, n))
		for i := 0; i < n; i++ {
			if err := x.DecodeValue(v.Index(i).Addr()); err != nil {
				return err
			}
		}

	case k == reflect.String:
		// get length
		n, err := x.decodeInt()
		if err != nil {
			return err
		}

		// read bytes
		r = make([]byte, n)
		if _, err := x.src.Read(r); err != nil {
			return err
		}

		// set value
		v.SetString(string(r))

	case k == reflect.Map:
		// get length
		n, err := x.decodeInt()
		if err != nil {
			return err
		}

		// make new map to ensure there is no unwanted data
		v.Set(reflect.MakeMap(t))
		kt := t.Key()
		vt := t.Elem()

		// populate map
		for i := 0; i < n; i++ {
			key := reflect.New(kt)
			if err := x.DecodeValue(key); err != nil {
				return err
			}

			val := reflect.New(vt)
			if err := x.DecodeValue(val); err != nil {
				return err
			}

			v.SetMapIndex(key.Elem(), val.Elem())
		}

	case k == reflect.Struct:
		// get length
		n, err := x.decodeInt()
		if err != nil {
			return err
		}

		if n != v.NumField() {
			return errors.New("field count missmatch")
		}

		for i := 0; i < n; i++ {
			if err := x.DecodeValue(v.Field(i).Addr()); err != nil {
				return err
			}
		}

	default:
		return errors.New("unsupported type")
	}

	return nil
}
