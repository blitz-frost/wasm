package dom

import (
	"syscall/js"
)

// An Event wraps a JS event object
type Event struct {
	js.Value
}

func (x Event) Cancel() {
	x.Call("preventDefault")
}

func (x Event) Target() Element {
	return Element{x.Get("target")}
}

type KeyboardEvent struct {
	Event
}

func AsKeyboardEvent(e Event) KeyboardEvent {
	return KeyboardEvent{e}
}

func (x KeyboardEvent) Key() string {
	return x.Get("key").String()
}

// A Handler wraps a JS event handler function.
type Handler struct {
	f js.Func
}

func NewHandler(fn func(this Element, e Event)) Handler {
	return Handler{js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() { fn(Element{this}, Event{args[0]}) }()
		return nil
	})}
}
