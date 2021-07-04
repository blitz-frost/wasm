package dom

import (
	"syscall/js"
)

type EventName string

const (
	EventBlur      EventName = "blur"
	EventClick               = "click"
	EventFocusIn             = "focusin"
	EventFocusOut            = "focusout"
	EventInput               = "input"
	EventKeyDown             = "keydown"
	EventMouseDown           = "mousedown"
	EventMouseUp             = "mouseup"
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

func (x KeyboardEvent) Key() string {
	return x.Get("key").String()
}

// A Handler wraps a JS event handler function.
type Handler struct {
	f js.Func
}

func MakeHandler(fn func(this Element, e Event)) Handler {
	return Handler{js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() { fn(Element{this}, Event{args[0]}) }()
		return nil
	})}
}
