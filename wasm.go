// The wasm package wraps the standard syscall/js package, and aims to provide a baseline API for interacting with the JS runtime from Go.
package wasm

import (
	"sync"
	"time"
	"unsafe"

	"syscall/js"

	"github.com/blitz-frost/errors"
	"github.com/blitz-frost/io"
	"github.com/blitz-frost/resource"
)

var Global = Object(js.Global())

// Any is used for documentation purposes, to clarify that it refers to a type handled by [js.ValueOf].
type Any = any

var ArrayClass = ClassGet("Array")

type Array Value

func ArrayAssert(v Value) (Array, bool) {
	ok := ArrayClass.CallRaw("isArray", v).Bool()
	return Array(v), ok
}

// length must be a value compatible with uint32
func ArrayMake(length int) Array {
	// there is no point in trying to use the other constructor variant (that takes the array elements), because js.ValueOf([]any) uses this variant + Value.SetIndex() on every element anyway

	o := ArrayClass.NewRaw(length)
	return Array(o)
}

func (x Array) Index(i int) Value {
	return Value(x).Index(i)
}

func (x Array) IndexSet(i int, v Any) {
	Value(x).SetIndex(i, v)
}

func (x Array) Length() int {
	return Value(x).Length()
}

type Bool Value

func BoolMake(v bool) Bool {
	o := js.ValueOf(v)
	return Bool(o)
}

func (x Bool) Get() bool {
	return Value(x).Bool()
}

var bytesClass = ClassGet("Uint8Array")

// Bytes mimics []byte using a JS Uint8Array as the underlying array.
type Bytes struct {
	v   Value
	len int
	cap int
}

func BytesOf(b []byte) Bytes {
	x := BytesMake(len(b), cap(b))
	x.CopyFrom(b)
	return x
}

func BytesMake(length, capacity int) Bytes {
	v := bytesClass.NewRaw(capacity)
	return Bytes{Value(v), length, capacity}
}

func View(arrayBuffer Value) Bytes {
	v := Value(bytesClass.NewRaw(arrayBuffer))
	n := v.Length()
	return Bytes{v, n, n}
}

func (x Bytes) Append(b []byte) Bytes {
	length := len(b) + x.len
	if length <= x.cap {
		// have room in current array
		v := x.v.Call("subarray", x.len, length)
		js.CopyBytesToJS(v, b)
		x.len = length
		return x
	}

	// not enough room; allocate new array and copy everything into it
	v := bytesClass.NewRaw(length)
	v.CallRaw("set", x.v)

	sub := v.CallRaw("subarray", x.len)
	js.CopyBytesToJS(sub, b)

	return Bytes{Value(v), length, length}
}

func (x Bytes) Cap() int {
	return x.cap
}

func (x Bytes) CopyFrom(b []byte) int {
	if len(b) > x.len {
		b = b[:x.len]
	}
	return js.CopyBytesToJS(x.v, b)
}

func (x Bytes) CopyTo(b []byte) int {
	if len(b) > x.len {
		b = b[:x.len]
	}
	return js.CopyBytesToGo(b, x.v)
}

func (x Bytes) Len() int {
	return x.len
}

func (x Bytes) Slice(start, end int) Bytes {
	v := x.v.Call("subarray", start)
	return Bytes{v, end - start, x.cap - start}
}

func (x Bytes) Value() Value {
	return x.v.Call("subarray", 0, x.len)
}

// BytesReader wraps a Bytes object to function as an [io.Reader].
// [Src] must be a valid Bytes value. It can be retrieved or exchanged when done, and will always be the remaining subslice of the initial data.
type BytesReader struct {
	Src Bytes
}

func (x *BytesReader) Read(b []byte) *io.Error {
	n := x.Src.CopyTo(b)
	x.Src = x.Src.Slice(n, x.Src.Len())
	if n < len(b) {
		return &io.Error{n, io.EOF}
	}
	return nil
}

// BytesWriter wraps a Bytes object to function as an [io.Writer].
// [Dst] must be a valid Bytes value. It may be freely retrieved or exchanged when done writing to the current value.
type BytesWriter struct {
	Dst Bytes
}

func (x *BytesWriter) Write(b []byte) *io.Error {
	x.Dst = x.Dst.Append(b)
	return nil
}

var catchNew = Function(Global.Get("goCatchNew"))

type Class Value

func ClassGet(class string) Class {
	return Class(Global.Get(class))
}

func (x Class) Call(method string, args ...Any) (Value, error) {
	return Object(x).Call(method, args...)
}

func (x Class) CallRaw(method string, args ...Any) Value {
	return Object(x).CallRaw(method, args...)
}

func (x Class) New(args ...Any) (Object, error) {
	r := catchNew.InvokeRaw(Value(x), args)
	o, err := catch(r)
	return Object(o), err
}

func (x Class) NewRaw(args ...Any) Object {
	o := Value(x).New(args...)
	return Object(o)
}

var dynamic = Export(func(this Value, args []Value) (Any, error) {
	p := resource.Pointer(args[0].Int())
	fn := p.Get().(Executor)
	return fn(this, args[1:])
})

// A DynamicFunction uses an indirect calling mechanism to execute an Executor.
// DynamicFunctions use a resource allocation system that can be infinitely reused, so can be used to safely create instantiated closures.
//
// The zero value is valid, but not useful until remade.
// Must be wiped when no longer needed.
type DynamicFunction struct {
	Function
	_p resource.Pointer
}

// fn may be nil, which will return the zero value.
func DynamicFunctionMake(fn Executor) DynamicFunction {
	if fn == nil {
		return DynamicFunction{}
	}

	p := resource.Alloc(fn)
	return DynamicFunction{
		Function: ClosureMake(dynamic.Function, uint(p)),
		_p:       p,
	}
}

// Remake wipes the old resource and recreates the DynamicFunction to use a new Executor.
func (x *DynamicFunction) Remake(fn Executor) {
	x.Wipe()
	*x = DynamicFunctionMake(fn)
}

func (x *DynamicFunction) Wipe() {
	if Value(x.Function).IsUndefined() {
		return
	}

	x._p.Wipe()
}

var (
	export            = Function(Global.Get("goExport"))
	exportNext uint32 = 1
	exportMux  sync.Mutex
)

// An ExportedFunction is a JS function that calls an Executor. It must be wiped when no longer needed.
//
// A program can create at most 2^32-1 ExportedFunctions. Because this is a finite resource, only unique functions should be exported. For instantiated closures, use [DynamicFunction].
type ExportedFunction struct {
	Function
	_f js.Func
}

// Export creates a Function that calls an Executor. If the Executor returns an error, it is thrown to the JS side (and the return value is discarded).
//
// See [js.FuncOf] for details regarding the relationship between the Go and JS runtimes during a Function call.
func Export(fn Executor) ExportedFunction {
	// base JS wrapper
	wrap := func(this Value, args []Value) Any {
		o, err := fn(this, args)
		if err != nil {
			jsErr := jsError(err)
			return []any{jsErr, true}
		}
		return []Any{o, false}
	}

	o := ExportRaw(wrap)

	// additional wrapper that throws on error
	v := export.InvokeRaw(Value(o.Function))
	o.Function = Function(v)

	return o
}

// ExportRaw wraps [js.FuncOf]. Unlike [Export], it does not add error handling layers, making the resulting Function slighly more efficient.
//
// js.FuncOf allocates internal IDs in a simple sequential manner, which can overflow.
// This can potentially overwrite old functions, which can result in undefined behaviour when calling or releasing them.
// This function allows code to safeguard against this by panicking when called too many times.
// Programs should be mindful of the amount of exported functions created. In most cases this should not really be a problem, but it's good to be aware that this resource is not infinite.
// Relevant as of Go 1.23.
func ExportRaw(fn func(Value, []Value) Any) ExportedFunction {
	exportMux.Lock()
	if exportNext == 0 {
		exportMux.Unlock()
		panic("JS function overflow")
	}
	exportNext++
	exportMux.Unlock()

	var o ExportedFunction
	o._f = js.FuncOf(fn)
	o.Function = Function(o._f.Value)
	return o
}

func (x ExportedFunction) Wipe() {
	x._f.Release()
}

var functionInvokeCatch = Function(Global.Get("goCatchInvoke"))

// Function is a JS function.
// In documentation, "Function" can be used in the same way as "func" to describe expected function signatures.
// Note that JS functions return at most one value.
type Function Value

var (
	asyncClosure = Function(Global.Get("goAsyncClosure"))
	closure      = Function(Global.Get("goClosure"))
)

// AsyncMake returns a Function(...args) => goAsync(fn, ...args)
//
// fn - Function(resolve, reject Function(Value), args ...Value)
func AsyncMake(fn Function) Function {
	o := asyncClosure.InvokeRaw(Value(fn))
	return Function(o)
}

// ClosureMake returns a Function(...args) => fn(data, ...args)
//
// fn - Function(Value, ...Value)
func ClosureMake(fn Function, data Any) Function {
	o := closure.InvokeRaw(Value(fn), data)
	return Function(o)
}

func (x Function) Invoke(args ...Any) (Value, error) {
	r := functionInvokeCatch.InvokeRaw(Value(x), args)
	return catch(r)
}

func (x Function) InvokeRaw(args ...Any) Value {
	return Value(x).Invoke(args...)
}

func (x Function) Value() Value {
	return Value(x)
}

// An Indirect can be used to create reusable ExportedFunctions, for a small performance cost.
//
// Its methods are concurrent safe.
type Indirect struct {
	_fn  Executor
	_mux sync.Mutex
}

func IndirectMake(fn Executor) Indirect {
	return Indirect{
		_fn: fn,
	}
}

func (x *Indirect) Exec(this Value, args []Value) (Any, error) {
	x._mux.Lock()
	fn := x._fn
	x._mux.Unlock()
	return fn(this, args)
}

// Set replaces the current underlying Executor.
func (x *Indirect) Set(fn Executor) {
	x._mux.Lock()
	x._fn = fn
	x._mux.Unlock()
}

type Executor func(Value, []Value) (Any, error)

// AsyncExec wraps an Executor to be used with "goAsync", to execute Go code asynchronously.
// The given Executor will be called in its own goroutine, and its return values will be used to resolve or reject the corresponding promise.
//
// Returns a Function(Function(Value), Function(Value), ...Value), suitable for [AsyncMake].
func AsyncExecutor(fn Executor) Executor {
	return func(this Value, args []Value) (Any, error) {
		go func() {
			o, err := fn(this, args[2:])
			if err == nil {
				resolve := args[0]
				resolve.Invoke(o)
			} else {
				reject := args[1]
				jsErr := jsError(err)
				reject.Invoke(jsErr)
			}
		}()
		return nil, nil
	}
}

// The underlying representation is float64.
type Number Value

func NumberMake[T NumberConstraint](v T) Number {
	o := js.ValueOf(v)
	return Number(o)
}

func (x Number) Float() float64 {
	return Value(x).Float()
}

func (x Number) Int() int {
	return Value(x).Int()
}

// js.ValueOf() uses a type switch and will not recognize defined types
type NumberConstraint interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | unsafe.Pointer | float32 | float64
}

var (
	object          = ClassGet("Object")
	objectCallCatch = Function(Global.Get("goCatchCall"))
)

type Object Value

func ObjectMake() Object {
	// this is what js.ValueOf() does, followed by setting the fields one-by-one
	return object.NewRaw()
}

func (x Object) Call(method string, args ...Any) (Value, error) {
	r := objectCallCatch.InvokeRaw(Value(x), method, args)
	return catch(r)
}

func (x Object) CallRaw(method string, args ...Any) Value {
	return Value(x).Call(method, args...)
}

func (x Object) Delete(key string) {
	Value(x).Delete(key)
}

func (x Object) Get(key string) Value {
	return Value(x).Get(key)
}

func (x Object) Is(class Class) bool {
	return Value(x).InstanceOf(Value(class))
}

// Keys returns the keys of a JS object.
func (x Object) Keys() []string {
	keys := object.CallRaw("keys", Value(x))
	n := keys.Length()
	o := make([]string, n)
	for i := range n {
		o[i] = keys.Index(i).String()
	}

	return o
}

func (x Object) Set(key string, value Any) {
	Value(x).Set(key, value)
}

var promise = ClassGet("Promise")

// In documentation "Promise(T)" can be used to describe a Promise that resolves to a value of type T.
type Promise Value

// PromiseMake returns a new promise.
// exec - Function(resolve, reject Function(Value))
func PromiseMake(exec Function) (Promise, error) {
	o, err := promise.New(Value(exec))
	return Promise(o), err
}

// Promise resolve returns a promise that immediately resolves with the given value.
func PromiseResolve(result Any) Promise {
	o := promise.CallRaw("resolve", result)
	return Promise(o)
}

var promiseAwait = Export(func(this Value, args []Value) (Any, error) {
	data := args[0]
	resolved := data.Index(0).Bool()
	p := PointerFrom(Number(data.Index(1)))

	ch := p.Get().(chan promiseAwaitRes)
	ch <- promiseAwaitRes{
		resolved: resolved,
		value:    args[1],
	}

	return Value{}, nil
})

type promiseAwaitRes struct {
	resolved bool
	value    Value
}

// Await synchronizes the Promise with the calling goroutine.
//
// Must not be called from the JS event loop.
func (x Promise) Await() (Value, error) {
	ch := make(chan promiseAwaitRes, 1)
	p := resource.Alloc(ch)

	resolveData := []any{true, uint(p)}
	resolve := ClosureMake(promiseAwait.Function, resolveData)

	rejectData := []any{false, uint(p)}
	reject := ClosureMake(promiseAwait.Function, rejectData)

	x.ThenFull(resolve, reject)
	res := <-ch
	p.Wipe()

	// both resolve and reject can be undefined values
	if res.resolved {
		return res.value, nil
	}

	var msg string
	if res.value.IsUndefined() {
		msg = "unspecified"
	} else {
		msg = res.value.Get("message").String()
	}
	return Value{}, errors.Simple(msg)
}

// fulfilled - Function(Value)
func (x Promise) Then(fulfilled Function) Promise {
	o := Value(x).Call("then", Value(fulfilled))
	return Promise(o)
}

// fulfilled, rejected - Function(Value)
func (x Promise) ThenFull(fulfilled, rejected Function) Promise {
	o := Value(x).Call("then", Value(fulfilled), Value(rejected))
	return Promise(o)
}

type String Value

func StringMake(v string) String {
	return String(js.ValueOf(v))
}

func (x String) Get() string {
	return Value(x).String()
}

// A Ticker represents a JS Interval. Useful to synchronize with the JS event loop.
type Ticker struct {
	V Value
}

// fn will be called with the provided args
func TickerMake(d time.Duration, fn Function, args ...Any) Ticker {
	var o Ticker
	ms := d.Milliseconds()
	callArgs := make([]any, 0, 2+len(args))
	callArgs = append(callArgs, Value(fn), ms)
	callArgs = append(callArgs, args...)
	o.V = Global.CallRaw("setInterval", callArgs...)
	return o
}

// Stop disables the Ticker. Does not guarantee that the Ticker will not fire one more time.
// The next call to [EventLoopWait] ensures that the Ticker has fired for the last time.
// The registered Function must handle any desired explicit synchronization.
//
// NOTE It might be guaranteed when called from the event loop.
func (x Ticker) Stop() {
	Global.CallRaw("clearInterval", x.V)
}

// A Timer represents a JS Timeout. Useful to synchronize with the JS event loop.
type Timer struct {
	V Value
}

// fn will be called with the provided args
func TimerMake(d time.Duration, fn Function, args ...Any) Timer {
	var o Timer
	ms := d.Milliseconds()
	callArgs := make([]any, 0, 2+len(args))
	callArgs = append(callArgs, Value(fn), ms)
	callArgs = append(callArgs, args...)
	o.V = Global.CallRaw("setTimeout", callArgs...)
	return o
}

// Stop prevents the Timer from firing, if it has not already done so.
// The next call to [EventLoopWait] ensures that the Timer has either fired or will never fire.
// The registered Function must handle any desired explicit synchronization.
//
// NOTE It might be guaranteed when called from the event loop.
func (x Timer) Stop() {
	Global.CallRaw("clearTimeout", x.V)
}

type Type = js.Type

type Value = js.Value

func Copy(dst Bytes, src Bytes) {
	// clip overflow
	if src.len > dst.len {
		src.len = dst.len
	}

	v := src.v.Call("subarray", 0, src.len)

	dst.v.Call("set", v)
}

var eventLoopWait = Export(func(this Value, args []Value) (Any, error) {
	p := resource.Pointer(args[0].Int())
	ch := p.Get().(chan struct{})
	ch <- struct{}{}
	return Value{}, nil
})

// EventLoopWait returns as soon as the JS event microqueue is available.
// Useful for some synchronization operations.
//
// Must not be called from the JS event loop.
func EventLoopWait() {
	ch := make(chan struct{}, 1)
	p := resource.Alloc(ch)

	promise := PromiseResolve(uint(p))
	promise.Then(eventLoopWait.Function)

	<-ch
	p.Wipe()
}

// LoadDone must be called when the Go code is ready for use by the JS side.
//
// Note that the wasm module terminates when the main goroutine exits.
func LoadDone() {
	Global.CallRaw("goLoadDone")
}

// The argument must have been obtained through [PointerValue].
func PointerFrom(v Number) resource.Pointer {
	return resource.Pointer(v.Int())
}

func PointerValue(p resource.Pointer) Number {
	return NumberMake(uint(p))
}

var console = Object(Global.Get("console"))

// Print uses the console.log function to print JS values.
func Print(v Value) {
	console.CallRaw("log", v)
}

func catch(result Value) (Value, error) {
	if result.Index(0).Bool() {
		return Value{}, js.Error{result.Index(1)}
	}
	return result.Index(1), nil
}

var jsErrors = ClassGet("Error")

func jsError(err error) Value {
	o := jsErrors.NewRaw(err.Error())
	return Value(o)
}
