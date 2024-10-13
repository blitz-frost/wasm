// Package svg wraps SVG DOM elements.
package svg

import (
	"strconv"

	"github.com/blitz-frost/wasm"
	"github.com/blitz-frost/wasm/css"
)

const xmlns = "http://www.w3.org/2000/svg"

var doc = wasm.Object(wasm.Global.Get("document"))

func fmtLength(val uint16, unit css.Unit) string {
	return strconv.FormatUint(uint64(val), 10) + string(unit)
}

type Element wasm.Value

type Line Element

func MakeLine() Line {
	v := doc.CallRaw("createElementNS", xmlns, "line")
	return Line(v)
}

func (x Line) X0(val uint16, unit css.Unit) {
	wasm.Value(x).Call("setAttribute", "x1", fmtLength(val, unit))
}

func (x Line) X1(val uint16, unit css.Unit) {
	wasm.Value(x).Call("setAttribute", "x2", fmtLength(val, unit))
}

func (x Line) Y0(val uint16, unit css.Unit) {
	wasm.Value(x).Call("setAttribute", "y1", fmtLength(val, unit))
}

func (x Line) Y1(val uint16, unit css.Unit) {
	wasm.Value(x).Call("setAttribute", "y2", fmtLength(val, unit))
}

func (x Line) Color(color css.Color) {
	wasm.Value(x).Call("setAttribute", "stroke", string(color))
}

func (x Line) Width(val uint16, unit css.Unit) {
	wasm.Value(x).Call("setAttribute", "stroke-width", fmtLength(val, unit))
}
