package dom

import "github.com/blitz-frost/wasm"

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

// TextInsert inserts the given string at the current cursor position.
func TextInsert(str string) {
	doc.Call("execCommand", "insertText", false, str)
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
