package wasm

import (
	"errors"
	"sync"
	"time"

	"syscall/js"

	"github.com/blitz-frost/io"
	"github.com/blitz-frost/resource"
)

var Global = Object(js.Global())

var (
	array   = Global.Get("Uint8Array")
	console = Global.Get("console")
)

// Any is used for documentation purposes, to clarify that it refers to a type handled by [js.ValueOf].
type Any = any

// An AsyncInterface can be used with "goAsync" to execute Go code asynchronously.
// The underlying Interface will be executed in its own goroutine, and its return values will be used to resolve or reject the corresponding promise.
type AsyncInterface struct {
	V Interface
}

// Function(Function(Value), Function(Value), ...Value)
func (x AsyncInterface) Exec(this Value, args []Value) (Any, error) {
	go func(this Value, args []Value) {
		o, err := x.V.Exec(this, args[2:])
		if err == nil {
			resolve := args[0]
			resolve.Invoke(o)
		} else {
			reject := args[1]
			jsErr := jsError(err)
			reject.Invoke(jsErr)
		}
	}(this, args)
	return nil, nil
}

// Bytes mimics []byte using a JS Uint8Array as the underlying array.
type Bytes struct {
	v        Value
	length   int
	capacity int
}

func BytesOf(b []byte) Bytes {
	x := BytesMake(len(b), cap(b))
	x.CopyFrom(b)
	return x
}

func BytesMake(length, capacity int) Bytes {
	v := array.New(capacity)
	return Bytes{v, length, capacity}
}

func View(arrayBuffer Value) Bytes {
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

func (x Bytes) Len() int {
	return x.length
}

func (x Bytes) Slice(start, end int) Bytes {
	v := x.v.Call("subarray", start)
	return Bytes{v, end - start, x.capacity - start}
}

func (x Bytes) Value() Value {
	return x.v.Call("subarray", 0, x.length)
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

var dynamic = Export(InterfaceFunc(dynamicExec))

// A DynamicFunction uses an indirect calling mechanism to execute an Interface.
// DynamicFunctions use a resource allocation system that can be infinitely reused, so can be used to safely create instantiated closures.
//
// The zero value is valid, but not useful until remade.
// Must be wiped when no longer needed.
type DynamicFunction struct {
	Function
	p resource.Pointer
}

// inter may be nil, which will return the zero value.
func DynamicFunctionMake(inter Interface) *DynamicFunction {
	if inter == nil {
		return &DynamicFunction{}
	}

	p := resource.Alloc(inter)
	return &DynamicFunction{
		Function: ClosureMake(dynamic.Function, uint(p)),
		p:        p,
	}
}

// Remake wipes the old resource and recreates the DynamicFunction to use a new Interface.
func (x *DynamicFunction) Remake(inter Interface) {
	x.Wipe()
	*x = *DynamicFunctionMake(inter)
}

func (x *DynamicFunction) Wipe() {
	if Value(x.Function).IsUndefined() {
		return
	}

	x.p.Wipe()
}

func dynamicExec(this Value, args []Value) (Any, error) {
	p := resource.Pointer(args[0].Int())
	inter := p.Get().(Interface)
	return inter.Exec(this, args[1:])
}

var (
	export            = Global.Get("goExport")
	exportNext uint32 = 1
	exportMux  sync.Mutex
)

// An ExportedFunction is a JS function that calls Go code. It must be wiped when no longer needed.
//
// A program can create at most 2^32-1 ExportedFunctions. Because this is a finite resource, only unique functions should be exported. For instantiated closures, use [DynamicFunction].
type ExportedFunction struct {
	Function
	f js.Func
}

// Export creates a Function that executes an Interface. If the Interface returns an error, it is thrown to the JS side (and the return value is discarded).
//
// See [js.FuncOf] for details regarding the relationship between the Go and JS runtimes during a Function call.
func Export(inter Interface) ExportedFunction {
	// base JS wrapper
	fn := func(this Value, args []Value) Any {
		o, err := inter.Exec(this, args)
		if err != nil {
			jsErr := jsError(err)
			return []any{jsErr, true}
		}
		return []any{o, false}
	}

	o := ExportRaw(fn)

	// additional wrapper that throws on error
	v := export.Invoke(o.Function.Value())
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
	o.f = js.FuncOf(fn)
	o.Function = Function(o.f.Value)
	return o
}

func (x ExportedFunction) Wipe() {
	x.f.Release()
}

var functionInvokeCatch = Global.Get("goCatchInvoke")

// Function is a JS function.
// In documentation, "Function" can be used in the same way as "func" to describe expected function signatures.
// Note that JS functions return at most one value.
type Function Value

var (
	asyncClosure = Global.Get("goAsyncClosure")
	closure      = Global.Get("goClosure")
)

// AsyncMake returns a Function(...args) => goAsync(fn, ...args)
//
// fn - Function(resolve, reject Function(Value), args ...Value)
func AsyncMake(fn Function) Function {
	o := asyncClosure.Invoke(Value(fn))
	return Function(o)
}

// ClosureMake returns a Function(...args) => fn(data, ...args)
//
// fn - Function(Value, ...Value)
func ClosureMake(fn Function, data Any) Function {
	o := closure.Invoke(Value(fn), data)
	return Function(o)
}

func (x Function) Invoke(args ...Any) (Value, error) {
	r := functionInvokeCatch.Invoke(Value(x), args)
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
	inter Interface
	mux   sync.Mutex
}

func IndirectMake(inter Interface) *Indirect {
	return &Indirect{
		inter: inter,
	}
}

func (x *Indirect) Exec(this Value, args []Value) (Any, error) {
	x.mux.Lock()
	inter := x.inter
	x.mux.Unlock()
	return inter.Exec(this, args)
}

// Set replaces the current underlying Interface.
func (x *Indirect) Set(inter Interface) {
	x.mux.Lock()
	x.inter = inter
	x.mux.Unlock()
}

// An Interface can bridge between a JS function call and Go code.
type Interface interface {
	Exec(this Value, args []Value) (Any, error)
}

type InterfaceFunc func(Value, []Value) (Any, error)

func (x InterfaceFunc) Exec(this Value, args []Value) (Any, error) {
	return x(this, args)
}

var (
	object          = Global.Get("Object")
	objectCallCatch = Global.Get("goCatchCall")
)

type Object Value

func ObjectMake(fields map[string]Any) Object {
	return Object(js.ValueOf(fields))
}

func (x Object) Call(method string, args ...Any) (Value, error) {
	r := objectCallCatch.Invoke(Value(x), method, args)
	return catch(r)
}

func (x Object) CallRaw(method string, args ...Any) Value {
	return Value(x).Call(method, args...)
}

func (x Object) Get(key string) Value {
	return Value(x).Get(key)
}

// Keys returns the keys of a JS object.
func (x Object) Keys() []string {
	keys := object.Call("keys", Value(x))
	n := keys.Length()
	o := make([]string, n)
	for i := 0; i < n; i++ {
		o[i] = keys.Index(i).String()
	}

	return o
}

func (x Object) Set(key string, value Any) {
	Value(x).Set(key, value)
}

func (x Object) Value() Value {
	return Value(x)
}

var promise = Global.Get("Promise")

// In documentation "Promise(T)" can be used to describe a Promise that resolves to a value of type T.
type Promise Value

// PromiseMake returns a new promise.
// exec - Function(resolve, reject Function(Value))
func PromiseMake(exec Function) (Promise, error) {
	o, err := New(promise, Value(exec))
	return Promise(o), err
}

func PromiseResolve(value Any) Promise {
	o := promise.Call("resolve", value)
	return Promise(o)
}

var promiseAwaitExport = Export(InterfaceFunc(promiseAwait))

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
	resolve := ClosureMake(promiseAwaitExport.Function, resolveData)

	rejectData := []any{false, uint(p)}
	reject := ClosureMake(promiseAwaitExport.Function, rejectData)

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
	return Value{}, errors.New(msg)
}

func promiseAwait(this Value, args []Value) (Any, error) {
	data := args[0]
	resolved := data.Index(0).Bool()
	p := PointerFrom(data.Index(1))

	ch := p.Get().(chan promiseAwaitRes)
	ch <- promiseAwaitRes{
		resolved: resolved,
		value:    args[1],
	}

	return Value{}, nil
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

func (x Promise) Value() Value {
	return Value(x)
}

// A Ticker represents a JS Interval. Useful to synchronize with the JS event loop.
type Ticker struct {
	v Value
}

// fn will be called with the provided args
func TickerMake(d time.Duration, fn Function, args ...Any) Ticker {
	var o Ticker
	ms := d.Milliseconds()
	callArgs := make([]any, 0, 2+len(args))
	callArgs = append(callArgs, Value(fn), ms)
	callArgs = append(callArgs, args...)
	o.v = Global.CallRaw("setInterval", callArgs...)
	return o
}

// Stop disables the Ticker. Does not guarantee that the Ticker will not fire one more time.
// The next call to [EventLoopWait] ensures that the Ticker has fired for the last time.
// The registered Function must handle any desired explicit synchronization.
//
// NOTE It might be guaranteed when called from the event loop.
func (x Ticker) Stop() {
	Global.CallRaw("clearInterval", x.v)
}

// A Timer represents a JS Timeout. Useful to synchronize with the JS event loop.
type Timer struct {
	v    Value
	f    js.Func
	done bool
}

// fn will be called with the provided args
func TimerMake(d time.Duration, fn Function, args ...Any) Timer {
	var o Timer
	ms := d.Milliseconds()
	callArgs := make([]any, 0, 2+len(args))
	callArgs = append(callArgs, Value(fn), ms)
	callArgs = append(callArgs, args...)
	o.v = Global.CallRaw("setTimeout", callArgs...)
	return o
}

// Stop prevents the Timer from firing, if it has not already done so.
// The next call to [EventLoopWait] ensures that the Timer has either fired or will never fire.
// The registered Function must handle any desired explicit synchronization.
//
// NOTE It might be guaranteed when called from the event loop.
func (x Timer) Stop() {
	Global.CallRaw("clearTimeout", x.v)
}

type Value = js.Value

func Copy(dst Bytes, src Bytes) {
	// clip overflow
	if src.length > dst.length {
		src.length = dst.length
	}

	v := src.v.Call("subarray", 0, src.length)

	dst.v.Call("set", v)
}

var eventLoopWaitExport = Export(InterfaceFunc(eventLoopWait))

// EventLoopWait returns as soon as the JS event microqueue is available.
// Useful for some synchronization operations.
//
// Must not be called from the JS event loop.
func EventLoopWait() {
	ch := make(chan struct{}, 1)
	p := resource.Alloc(ch)

	promise := PromiseResolve(uint(p))
	promise.Then(eventLoopWaitExport.Function)

	<-ch
	p.Wipe()
}

func eventLoopWait(this Value, args []Value) (Any, error) {
	p := resource.Pointer(args[0].Int())
	ch := p.Get().(chan struct{})
	ch <- struct{}{}
	return Value{}, nil
}

// LoadDone must be called when the Go code is ready for use by the JS side.
//
// Note that the wasm module terminates when the main goroutine exits.
func LoadDone() {
	Global.CallRaw("goLoadDone")
}

var catchNew = Global.Get("goCatchNew")

func New(class Value, args ...Any) (Value, error) {
	r := catchNew.Invoke(class, args)
	return catch(r)
}

// Print uses the console.log function to print JS values.
func Print(v Value) {
	console.Call("log", v)
}

func PointerFrom(v Value) resource.Pointer {
	return resource.Pointer(v.Int())
}

func PointerValue(p resource.Pointer) Value {
	return js.ValueOf(uint(p))
}

func catch(v Value) (Value, error) {
	if v.Index(0).Bool() {
		return Value{}, errorFrom(v.Index(1))
	}

	return v.Index(1), nil
}

func errorFrom(v Value) error {
	errStr := v.Get("name").String()
	errStr += ": " + v.Get("message").String()

	return errors.New(errStr)
}

var jsErrors = Global.Get("Error")

func jsError(err error) Value {
	return jsErrors.New(err.Error())
}
