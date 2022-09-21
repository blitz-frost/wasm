package wasm

import (
	"errors"

	"github.com/blitz-frost/io"
	"syscall/js"
)

var (
	global = js.Global()

	array       = global.Get("Uint8Array")
	console     = global.Get("console")
	catchCall   = global.Get("catchCall")
	catchInvoke = global.Get("catchInvoke")
	object      = global.Get("Object")
)

// Bytes mimics []byte using a JS Uint8Array as the underlying array.
type Bytes struct {
	v        js.Value
	length   int
	capacity int
}

func BytesOf(b []byte) Bytes {
	x := MakeBytes(len(b), cap(b))
	x.CopyFrom(b)
	return x
}

func MakeBytes(length, capacity int) Bytes {
	v := array.New(capacity)
	return Bytes{v, length, capacity}
}

func View(arrayBuffer js.Value) Bytes {
	v := array.New(arrayBuffer)
	n := v.Length()
	return Bytes{v, n, n}
}

func (x Bytes) Append(b []byte) Bytes {
	length := len(b) + x.length
	if length <= x.capacity {
		// have room in current array
		v := x.v.Call("subarray", x.length, length)
		js.CopyBytesToJS(v, b)
		x.length = length
		return x
	}

	// not enough room; allocate new array and copy everything into it
	v := array.New(length)
	v.Call("set", x.v)

	sub := v.Call("subarray", x.length)
	js.CopyBytesToJS(sub, b)

	return Bytes{v, length, length}
}

func (x Bytes) Cap() int {
	return x.capacity
}

func (x Bytes) CopyFrom(b []byte) int {
	if len(b) > x.length {
		b = b[:x.length]
	}
	return js.CopyBytesToJS(x.v, b)
}

func (x Bytes) CopyTo(b []byte) int {
	if len(b) > x.length {
		b = b[:x.length]
	}
	return js.CopyBytesToGo(b, x.v)
}

func (x Bytes) Js() js.Value {
	return x.v.Call("subarray", 0, x.length)
}

func (x Bytes) Len() int {
	return x.length
}

func (x Bytes) Slice(start, end int) Bytes {
	v := x.v.Call("subarray", start)
	return Bytes{v, end - start, x.capacity - start}
}

// BytesReader wraps a Bytes object to function as an [io.Reader].
// [Src] must be a valid Bytes value. It can be retrieved or exchanged when done, and will always be the remaining subslice of the initial data.
type BytesReader struct {
	Src Bytes
}

func (x *BytesReader) Close() error {
	return nil
}

func (x *BytesReader) Read(b []byte) (int, error) {
	n := x.Src.CopyTo(b)
	x.Src = x.Src.Slice(n, x.Src.Len())
	if n < len(b) {
		return n, io.EOF
	}
	return n, nil
}

// BytesWriter wraps a Bytes object to function as an [io.Writer].
// [Dst] must be a valid Bytes value. It may be freely retrieved or exchanged when done writing to the current value.
type BytesWriter struct {
	Dst Bytes
}

func (x *BytesWriter) Close() error {
	return nil
}

func (x *BytesWriter) Write(b []byte) (int, error) {
	x.Dst = x.Dst.Append(b)
	return len(b), nil
}

// A Ticker represents a JS Interval. Useful to synchronize with the main JS thread.
type Ticker struct {
	v js.Value
	f js.Func
}

func MakeTicker(ms uint64, fn func()) Ticker {
	f := js.FuncOf(func(this js.Value, args []js.Value) any {
		fn()
		return nil
	})
	return Ticker{
		v: global.Call("setInterval", f, ms),
		f: f,
	}
}

func (x Ticker) Stop() {
	global.Call("clearInterval", x.v)
	x.f.Release()
}

// A Timer represents a JS Timeout. Useful to synchronize with the main JS thread.
type Timer struct {
	v js.Value
	f js.Func
}

func MakeTimer(ms uint64, fn func()) Timer {
	f := js.FuncOf(func(this js.Value, args []js.Value) any {
		fn()
		return nil
	})
	return Timer{
		v: global.Call("setTimeout", f, ms),
		f: f,
	}
}

func (x Timer) Stop() {
	global.Call("clearTimeout", x.v)
	x.f.Release()
}

// Await synchronizes the input promise.
func Await(promise js.Value) (js.Value, error) {
	resolveCh := make(chan js.Value)
	resolve := js.FuncOf(func(this js.Value, args []js.Value) any {
		var o js.Value
		if len(args) > 0 {
			o = args[0]
		}
		resolveCh <- o
		return nil
	})

	rejectCh := make(chan js.Value)
	reject := js.FuncOf(func(this js.Value, args []js.Value) any {
		// there should always be an error when rejecting... right?
		rejectCh <- args[0]
		return nil
	})

	promise.Call("then", resolve, reject)
	var o js.Value
	var err error
	select {
	case o = <-resolveCh:
	case o = <-rejectCh:
		msg := o.Get("message").String()
		err = errors.New(msg)
		o = js.Value{}
	}

	resolve.Release()
	reject.Release()
	return o, err
}

// Call is the method variant of Invoke.
func Call(obj js.Value, method string, args ...any) (js.Value, error) {
	r := catchCall.Invoke(obj, method, args)
	return catch(r)
}

func Copy(dst Bytes, src Bytes) {
	// clip overflow
	if src.length > dst.length {
		src.length = dst.length
	}

	v := src.v.Call("subarray", 0, src.length)

	dst.v.Call("set", v)
}

// Invoke exectues a function call, catching a thrown exception and returning it as a Go error.
func Invoke(fn js.Value, args ...any) (js.Value, error) {
	r := catchInvoke.Invoke(fn, args)
	return catch(r)
}

// Keys returns the keys of a JS object.
func Keys(obj js.Value) []string {
	if obj.Type() != js.TypeObject {
		return nil
	}

	keys := object.Call("keys", obj)
	n := keys.Length()
	o := make([]string, n)
	for i := 0; i < n; i++ {
		o[i] = keys.Index(i).String()
	}

	return o
}

// Print uses the console.log function to print JS values.
func Print(v js.Value) {
	console.Call("log", v)
}

func catch(v js.Value) (js.Value, error) {
	if v.Index(0).Bool() {
		return js.Undefined(), errorFrom(v.Index(1))
	}

	return v.Index(1), nil
}

func errorFrom(v js.Value) error {
	errStr := v.Get("name").String()
	errStr += ": " + v.Get("message").String()

	return errors.New(errStr)
}
