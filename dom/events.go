package dom

import (
	"syscall/js"
)

// An Event wraps a JS event object
type Event struct {
	js.Value
}

func (x Event) GetTarget() *Base {
	return newBase(x.Get("target"))
}

// A Handler wraps a JS event handler function.
type Handler struct {
	f js.Func
}

func NewHandler(fn func(this *Base, e Event)) *Handler {
	return &Handler{js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() { fn(newBase(this), Event{args[0]}) }()
		return nil
	})}
}
