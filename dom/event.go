package dom

import (
	"github.com/blitz-frost/wasm"
)

type EventName string

const (
	EventBlur       EventName = "blur"
	EventChange               = "change"
	EventClick                = "click"
	EventClickRight           = "contextmenu"
	EventFocus                = "focus"
	EventFocusIn              = "focusin"
	EventFocusOut             = "focusout"
	EventInput                = "input"
	EventKeyDown              = "keydown"
	EventKeyUp                = "keyup"
	EventMouseDown            = "mousedown"
	EventMouseEnter           = "mouseenter"
	EventMouseLeave           = "mouseleave"
	EventMouseMove            = "mousemove"
	EventMouseUp              = "mouseup"
	EventMouseWheel           = "mousewheel"
	EventResize               = "resize"
)

// An Event wraps a JS event object
type Event struct {
	wasm.Object
}

func (x Event) Cancel() {
	x.CallRaw("stopPropagation")
}

func (x Event) CancelDefault() {
	x.CallRaw("preventDefault")
}

func (x Event) Target() Element {
	return Element{x.Get("target")}
}

type KeyboardEvent struct {
	Event
}

func (x KeyboardEvent) Code() string {
	return x.Get("code").String()
}

// Ctrl returns true if the Ctrl key is being pressed.
func (x KeyboardEvent) Ctrl() bool {
	return x.Get("ctrlKey").Bool()
}

func (x KeyboardEvent) Key() string {
	return x.Get("key").String()
}

type MouseEvent struct {
	Event
}

func (x MouseEvent) Button() byte {
	return byte(x.Get("button").Int())
}

func (x MouseEvent) XAbs() uint16 {
	return uint16(x.Get("pageX").Int())
}

func (x MouseEvent) YAbs() uint16 {
	return uint16(x.Get("pageY").Int())
}

func (x MouseEvent) XRel() uint16 {
	return uint16(x.Get("offsetX").Int())
}

func (x MouseEvent) YRel() uint16 {
	return uint16(x.Get("offsetY").Int())
}

type WheelEvent struct {
	Event
}

func (x WheelEvent) Y() int8 {
	return int8(x.Get("deltaY").Float())
}

type Handler interface {
	Handle(Event)
}

type HandlerFunc func(Event)

func (x HandlerFunc) Exec(this wasm.Value, args []wasm.Value) (wasm.Any, error) {
	e := Event{wasm.Object(args[0])}
	x(e)
	return nil, nil
}

func (x HandlerFunc) Handle(e Event) {
	x(e)
}

type HandlerInterface struct {
	Handler
}

func (x HandlerInterface) Exec(this wasm.Value, args []wasm.Value) (wasm.Any, error) {
	e := Event{wasm.Object(args[0])}
	x.Handle(e)
	return nil, nil
}

// A HandlerFunction is a Function(Event).
type HandlerFunction wasm.Function
