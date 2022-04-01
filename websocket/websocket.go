// Package websocket wraps the Javascript Websocket API.
package websocket

import (
	"fmt"
	"syscall/js"
	"time"
)

const (
	CodeNormal          Code = 1000
	CodeGoingAway            = 1001
	CodeProtocolError        = 1002
	CodeUnsupported          = 1003
	CodeReserved             = 1004
	CodeNoStatus             = 1005
	CodeAbnormal             = 1006
	CodeInvalidFrame         = 1007
	CodePolicyViolation      = 1008
	CodeTooBig               = 1009
	CodeMandatory            = 1010
	CodeInternalError        = 1011
	CodeRestart              = 1012
	CodeTryAgain             = 1013
	CodeBadGateway           = 1014
	CodeTLS                  = 1015
)

var (
	global = js.Global()
	alloc  = global.Get("Uint8Array")
)

type Code int

type Conn struct {
	v js.Value

	closeHandler js.Func
	closeFunc    func(Code)

	msgHandler js.Func
	binaryFunc func([]byte)
	textFunc   func(string)
}

func (x Conn) Close() error {
	x.v.Call("close")
	return nil
}

// HandleCloses registers fn to be called with the closing code, when the websocket closes.
func (x *Conn) OnClose(fn func(Code)) {
	x.closeFunc = fn
}

// HandleBinary registers fn to be called with incoming binary data.
func (x *Conn) OnBinary(fn func([]byte)) {
	x.binaryFunc = fn
}

// HandleText registers fn to be called with incoming text data.
func (x *Conn) OnText(fn func(string)) {
	x.textFunc = fn
}

// Wait blocks until the currently buffered data is written out, or the websocket closes.
// Since the JS API doesn't have an event for this, the current buffer amount is checked every d time.
func (x Conn) Wait(d time.Duration) {
	for {
		if x.v.Get("readyState").Int() == 3 {
			return
		}

		if x.v.Get("bufferedAmount").Int() == 0 {
			return
		}

		time.Sleep(d)
	}
}

// Write buffers data to be written out to the other end of the websocket.
func (x Conn) Write(b []byte) error {
	buf := alloc.New(len(b))
	js.CopyBytesToJS(buf, b)

	x.v.Call("send", buf)

	return nil
}

func (x Conn) WriteText(b []byte) error {
	x.v.Call("send", string(b))

	return nil
}

func Dial(url string) *Conn {
	ws := global.Get("WebSocket").New(url)
	ws.Set("binaryType", "arraybuffer") // requires ones less async function call to read
	x := Conn{v: ws}

	// hook close event
	x.closeFunc = func(c Code) { fmt.Println("websocket closed", c) }
	x.closeHandler = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		code := Code(args[0].Get("code").Int())
		x.closeFunc(code)

		// cleanup
		x.msgHandler.Release()
		x.closeHandler.Release()

		return nil
	})
	ws.Set("onclose", x.closeHandler)

	// hook message event
	x.binaryFunc = func([]byte) {}
	x.textFunc = func(string) {}
	x.msgHandler = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		data := args[0].Get("data")
		if data.Type() == js.TypeString {
			x.textFunc(data.String())
		} else {
			buf := alloc.New(data)
			n := buf.Get("byteLength").Int()
			b := make([]byte, n)
			js.CopyBytesToGo(b, buf)
			x.binaryFunc(b)
		}

		return nil
	})
	ws.Set("onmessage", x.msgHandler)

	// wait for connection
	ch := make(chan struct{})
	fn := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ch <- struct{}{}
		return nil
	})
	ws.Set("onopen", fn)

	<-ch
	fn.Release()
	return &x
}
