package wasm

import (
	"errors"

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

// Bytes wraps a JS Uint8Array.
type Bytes struct {
	v js.Value
}

func BytesFrom(b []byte) Bytes {
	x := MakeBytes(len(b))
	x.CopyFrom(b)
	return x
}

func MakeBytes(n int) Bytes {
	v := array.New(n)
	return Bytes{v}
}

func View(arrayBuffer js.Value) Bytes {
	v := array.New(arrayBuffer)
	return Bytes{v}
}

func (x Bytes) CopyFrom(b []byte) {
	js.CopyBytesToJS(x.v, b)
}

func (x Bytes) CopyTo(b []byte) {
	js.CopyBytesToGo(b, x.v)
}

func (x Bytes) Js() js.Value {
	return x.v
}

func (x Bytes) Length() int {
	return x.v.Length()
}

func (x Bytes) Slice(start, end int) Bytes {
	v := x.v.Call("subarray", start, end)
	return Bytes{v}
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
