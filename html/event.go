package html

import (
	"github.com/blitz-frost/wasm"
)

type EventName string

const (
	EventBlur       EventName = "blur"
	EventChange     EventName = "change"
	EventClick      EventName = "click"
	EventClickRight EventName = "contextmenu"
	EventFocus      EventName = "focus"
	EventFocusIn    EventName = "focusin"
	EventFocusOut   EventName = "focusout"
	EventInput      EventName = "input"
	EventKeyDown    EventName = "keydown"
	EventKeyUp      EventName = "keyup"
	EventMouseDown  EventName = "mousedown"
	EventMouseEnter EventName = "mouseenter"
	EventMouseLeave EventName = "mouseleave"
	EventMouseMove  EventName = "mousemove"
	EventMouseUp    EventName = "mouseup"
	EventMouseWheel EventName = "mousewheel"
	EventResize     EventName = "resize"
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
	v := x.Get("target")
	return Element{wasm.Object(v)}
}

type Handler func(Event)

func (x Handler) Exec(this wasm.Value, args []wasm.Value) (wasm.Any, error) {
	e := Event{wasm.Object(args[0])}
	x(e)
	return nil, nil
}

// A HandlerFunction is a Function(Event).
type HandlerFunction wasm.Function

type KeyboardEvent struct {
	Event
}

func (x KeyboardEvent) Code() string {
	v := x.Get("code")
	return v.String()
}

// Ctrl returns true if the Ctrl key is being pressed.
func (x KeyboardEvent) Ctrl() bool {
	v := x.Get("ctrlKey")
	return v.Bool()
}

func (x KeyboardEvent) Key() string {
	v := x.Get("key")
	return v.String()
}

type MouseEvent struct {
	Event
}

func (x MouseEvent) Button() byte {
	v := x.Get("button")
	return byte(v.Int())
}

func (x MouseEvent) XAbs() uint16 {
	v := x.Get("pageX")
	return uint16(v.Int())
}

func (x MouseEvent) YAbs() uint16 {
	v := x.Get("pageY")
	return uint16(v.Int())
}

func (x MouseEvent) XRel() uint16 {
	v := x.Get("offsetX")
	return uint16(v.Int())
}

func (x MouseEvent) YRel() uint16 {
	v := x.Get("offsetY")
	return uint16(v.Int())
}

type WheelEvent struct {
	Event
}

func (x WheelEvent) Y() int8 {
	v := x.Get("deltaY")
	return int8(v.Float())
}
