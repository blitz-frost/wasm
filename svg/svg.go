// Package svg wraps SVG DOM elements.
package svg

import (
	"strconv"
	"syscall/js"

	"github.com/blitz-frost/wasm/css"
)

const xmlns = "http://www.w3.org/2000/svg"

var doc js.Value = js.Global().Get("document")

type Element interface {
	JSValue() js.Value
}

func fmtLength(val uint16, unit css.Unit) string {
	return strconv.FormatUint(uint64(val), 10) + string(unit)
}

type Line struct {
	Value js.Value
}

func MakeLine() Line {
	return Line{doc.Call("createElementNS", xmlns, "line")}
}

func (x Line) X0(val uint16, unit css.Unit) {
	x.Value.Call("setAttribute", "x1", fmtLength(val, unit))
}

func (x Line) X1(val uint16, unit css.Unit) {
	x.Value.Call("setAttribute", "x2", fmtLength(val, unit))
}

func (x Line) Y0(val uint16, unit css.Unit) {
	x.Value.Call("setAttribute", "y1", fmtLength(val, unit))
}

func (x Line) Y1(val uint16, unit css.Unit) {
	x.Value.Call("setAttribute", "y2", fmtLength(val, unit))
}

func (x Line) Color(color css.Color) {
	x.Value.Call("setAttribute", "stroke", string(color))
}

func (x Line) Width(val uint16, unit css.Unit) {
	x.Value.Call("setAttribute", "stroke-width", fmtLength(val, unit))
}

func (x Line) JSValue() js.Value {
	return x.Value
}
