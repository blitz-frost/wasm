package dom

import (
	"strconv"

	"github.com/blitz-frost/wasm/css"
)

func fmtLength(val uint16, unit css.Unit) string {
	return strconv.FormatUint(uint64(val), 10) + string(unit)
}

type Style struct {
	m map[string]string
}

func NewStyle() Style {
	return Style{make(map[string]string)}
}

func (x Style) Copy() Style {
	o := NewStyle()

	for k, v := range x.m {
		o.m[k] = v
	}

	return o
}

func (x Style) sideLong(base, name string, val string, sides ...css.Side) {
	for _, side := range sides {
		x.m[base+string(side)+name] = val
	}
}

func (x Style) side(name string, val string, sides ...css.Side) {
	for _, side := range sides {
		x.m[name+string(side)] = val
	}
}

func (x Style) BackgroundColor(color css.Color) {
	x.m["backgroundColor"] = string(color)
}

func (x Style) Border(width uint16, unit css.Unit, style css.Border, color css.Color, sides ...css.Side) {
	val := fmtLength(width, unit) + " " + string(style) + " " + string(color)
	x.side("border", val, sides...)
}

func (x Style) BorderCollapse(val bool) {
	var str string
	if val {
		str = "collapse"
	} else {
		str = "separate"
	}
	x.m["borderCollapse"] = str

}

func (x Style) BorderColor(color css.Color, sides ...css.Side) {
	x.sideLong("border", "Color", string(color), sides...)
}

func (x Style) BorderRadius(val uint16, unit css.Unit, corners ...css.Corner) {
	str := fmtLength(val, unit)
	for _, corner := range corners {
		x.m["border"+string(corner)+"Radius"] = str
	}
}

func (x Style) BorderStyle(style css.Border, sides ...css.Side) {
	x.sideLong("border", "Style", string(style), sides...)
}

func (x Style) BorderWidth(width uint16, unit css.Unit, sides ...css.Side) {
	x.sideLong("border", "Width", fmtLength(width, unit), sides...)
}

func (x Style) Cursor(val css.Cursor) {
	x.m["cursor"] = string(val)
}

func (x Style) Float(val css.Float) {
	x.m["cssFloat"] = string(val)
}

func (x Style) Height(val uint16, unit css.Unit) {
	x.m["height"] = fmtLength(val, unit)
}

func (x Style) HeightMax(val uint16, unit css.Unit) {
	x.m["maxHeight"] = fmtLength(val, unit)
}

func (x Style) HeightMin(val uint16, unit css.Unit) {
	x.m["minHeight"] = fmtLength(val, unit)
}

func (x Style) Margin(val uint16, unit css.Unit, sides ...css.Side) {
	x.side("margin", fmtLength(val, unit), sides...)
}

func (x Style) Padding(val uint16, unit css.Unit, sides ...css.Side) {
	x.side("padding", fmtLength(val, unit), sides...)
}

func (x Style) Resize(val css.Resize) {
	x.m["resize"] = string(val)
}

func (x Style) TextAlign(val css.Align) {
	x.m["textAlign"] = string(val)
}

func (x Style) Width(val uint16, unit css.Unit) {
	x.m["width"] = fmtLength(val, unit)
}

func (x Style) WidthMax(val uint16, unit css.Unit) {
	x.m["maxWidth"] = fmtLength(val, unit)
}

func (x Style) WidthMin(val uint16, unit css.Unit) {
	x.m["minWidth"] = fmtLength(val, unit)
}
