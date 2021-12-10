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

// InsertText inserts the given string at the current cursor position.
func InsertText(str string) {
	doc.Call("execCommand", "insertText", false, str)
}

// SelectText selects text inside the current active element.
func SelectText(start, end int) {
	sel := window.Call("getSelection")
	if sel.Get("rangeCount").Int() == 0 {
		return
	}

	rng := sel.Call("getRangeAt", 0)
	node := rng.Get("endContainer")
	rng.Call("setEnd", node, end)
	rng.Call("setStart", node, start)
}
