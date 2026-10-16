package window

import (
	"net/url"

	"github.com/blitz-frost/wasm"
	"github.com/blitz-frost/wasm/html"
)

var location = wasm.Object(wasm.Global.Get("location"))

// CaretMove moves caret position inside the current selection.
func CaretMove(pos int) {
	sel := wasm.Global.CallRaw("getSelection")
	if sel.Get("rangeCount").Int() == 0 {
		return
	}

	rng := sel.Call("getRangeAt", 0)
	node := rng.Get("endContainer")
	rng.Call("setEnd", node, pos)
	rng.Call("setStart", node, pos)
}

func Handle(event html.EventName, h html.HandlerFunction) {
	wasm.Global.CallRaw("addEventListener", string(event), wasm.Value(h))
}

func HandleRemove(event html.EventName, h html.HandlerFunction) {
	wasm.Global.CallRaw("removeEventListener", string(event), wasm.Value(h))
}

// Url returns the current navigation URL.
func Url() url.URL {
	s := location.Get("href").String()
	u, _ := url.Parse(s)
	return *u
}

// TextSelect selects text inside the current active element.
func TextSelect(start, end int) {
	sel := wasm.Global.CallRaw("getSelection")
	if sel.Get("rangeCount").Int() == 0 {
		return
	}

	rng := sel.Call("getRangeAt", 0)
	node := rng.Get("endContainer")
	rng.Call("setEnd", node, end)
	rng.Call("setStart", node, start)
}
