// Package svg wraps SVG DOM elements.
package svg

import (
	"strconv"

	"github.com/blitz-frost/wasm"
	"github.com/blitz-frost/wasm/css"
	"github.com/blitz-frost/wasm/html"
	"github.com/blitz-frost/wasm/html/document"
)

const xmlns = "http://www.w3.org/2000/svg"

func fmtLength(val uint16, unit css.Unit) string {
	return strconv.FormatUint(uint64(val), 10) + string(unit)
}

type Container struct {
	Element
}

func ContainerMake() Container {
	v := document.ElementCreateNS(xmlns, "svg")
	return Container{Element(v)}
}

func (x Container) Append(elem Element) {
	x.CallRaw("appendChild", elem.Value())
}

func (x Container) Children() html.Collection {
	//NOTE this is a HTML property
	v := x.Get("children")
	return html.Collection{wasm.Object(v)}
}

// NOTE This is likely both an XML element and a HTML element, which are very similar, but not identical.
// For the time being we don't care enough about XML to create its own package, or attempt to factor it with HTML.
type Element struct {
	wasm.Object
}

func (x Element) attributeSet(name, value string) {
	x.CallRaw("setAttribute", name, value)
}

type Line struct {
	Element
}

func LineMake() Line {
	v := document.ElementCreateNS(xmlns, "line")
	return Line{Element(v)}
}

func (x Line) X0(val uint16, unit css.Unit) {
	x.attributeSet("x1", fmtLength(val, unit))
}

func (x Line) X1(val uint16, unit css.Unit) {
	x.attributeSet("x2", fmtLength(val, unit))
}

func (x Line) Y0(val uint16, unit css.Unit) {
	x.attributeSet("y1", fmtLength(val, unit))
}

func (x Line) Y1(val uint16, unit css.Unit) {
	x.attributeSet("y2", fmtLength(val, unit))
}

func (x Line) Color(color css.Color) {
	x.attributeSet("stroke", string(color))
}

func (x Line) Width(val uint16, unit css.Unit) {
	x.attributeSet("stroke-width", fmtLength(val, unit))
}
