package dom

import (
	"syscall/js"
)

type EventName string

const (
	EventBlur       EventName = "blur"
	EventChange               = "change"
	EventClick                = "click"
	EventClickRight           = "contextmenu"
	EventFocusIn              = "focusin"
	EventFocusOut             = "focusout"
	EventInput                = "input"
	EventKeyDown              = "keydown"
	EventMouseDown            = "mousedown"
	EventMouseUp              = "mouseup"
)

// An Event wraps a JS event object
type Event struct {
	js.Value
}

func (x Event) Cancel() {
	x.Call("stopPropagation")
}

func (x Event) CancelDefault() {
	x.Call("preventDefault")
}

func (x Event) Target() Element {
	return Element{x.Get("target")}
}

type KeyboardEvent struct {
	Event
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

func (x MouseEvent) X() uint16 {
	return uint16(x.Get("pageX").Float())
}

func (x MouseEvent) Y() uint16 {
	return uint16(x.Get("pageY").Float())
}

// A Handler wraps a JS event handler function.
type Handler struct {
	f js.Func
}

// MakeHandler wraps a Go function to be used as a DOM event handler.
// fn must be non blocking, otherwise the application will deadlock.
// Notably, http requests block.
func MakeHandler(fn func(Event)) Handler {
	return Handler{js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fn(Event{args[0]})
		return nil
	})}
}

// Delete releases the underlying JS function.
func (x Handler) Delete() {
	x.f.Release()
}
