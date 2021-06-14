package dom

// CaretMove moves caret position inside the current selection.
func CaretMove(pos int) {
	sel := window.Call("getSelection")
	if sel.Get("rangeCount").Int() == 0 {
		return
	}

	rng := sel.Call("getRangeAt", 0)
	node := rng.Get("endContainer")
	rng.Call("setEnd", node, pos)
	rng.Call("setStart", node, pos)
}
