// Package css provides CSS Go definitions.
package css

import (
	"strconv"
)

type Align string

const (
	AlignLeft    Align = "left"
	AlignRight         = "right"
	AlignCenter        = "center"
	AlignJustify       = "justify"
)

type BorderStyleKind string

const (
	BorderNone   BorderStyleKind = "none"
	BorderHidden                 = "hidden"
	BorderDotted                 = "dotted"
	BorderDashed                 = "dashed"
	BorderSolid                  = "solid"
	BorderDouble                 = "double"
	BorderGroove                 = "groove"
	BorderRidge                  = "ridge"
	BorderInset                  = "inset"
	BorderOutset                 = "outset"
)

type Color string

const (
	Transparent    Color = "transparent"
	Black                = "black"
	CadetBlue            = "cadetblue"
	CornflowerBlue       = "cornflowerblue"
	Cyan                 = "cyan"
	Salmon               = "salmon"
	White                = "white"
	WhiteSmoke           = "whitesmoke"
)

type Corner string

const (
	BottomLeft  Corner = "BottomLeft"
	BottomRight        = "BottomRight"
	TopLeft            = "TopLeft"
	TopRight           = "TopRight"
)

var CornerAll []Corner = []Corner{BottomLeft, BottomRight, TopLeft, TopRight}

type CursorKind string

const (
	CursorAlias        CursorKind = "alias"
	CursorScroll                  = "all-scroll"
	CursorCell                    = "cell"
	CursorContext                 = "context-menu"
	CursorCopy                    = "copy"
	CursorCross                   = "crosshair"
	CursorDefault                 = "default"
	CursorHelp                    = "help"
	CursorMove                    = "move"
	CursorNodrop                  = "no-drop"
	CursorNone                    = "none"
	CursorInvalid                 = "not-allowed"
	CursorPointer                 = "pointer"
	CursorProgress                = "progress"
	CursorResizeCol               = "col-resize"
	CursorResizeRow               = "row-resize"
	CursorResizeE                 = "e-resize"
	CursorResizeEW                = "ew-resize"
	CursorResizeN                 = "n-resize"
	CursorResizeNE                = "ne-resize"
	CursorResizeNESW              = "nesw-resize"
	CursorResizeNS                = "ns-resize"
	CursorResizeNW                = "nw-resize"
	CursorResizeNWSE              = "nwse-resize"
	CursorResizeS                 = "s-resize"
	CursorResizeSE                = "se-resize"
	CursorResizeSW                = "sw-resize"
	CursorResizeW                 = "w-resize"
	CursorText                    = "text"
	CursorTextVertical            = "vertical-text"
	CursorWait                    = "wait"
	CursorZoomIn                  = "zoom-in"
	CursorZoomOut                 = "zoom-out"
)

type DisplayKind string

const (
	DisplayBlock       DisplayKind = "block"
	DisplayFlex                    = "flex"
	DisplayGrid                    = "grid"
	DisplayInline                  = "inline"
	DisplayInlineBlock             = "inline-block"
	DisplayInlineFlex              = "inline-flex"
	DisplayInlineGrid              = "inline-grid"
	DisplayInlineTable             = "inline-table"
	DisplayTable                   = "table"
)

type FloatKind string

const (
	FloatNone  FloatKind = "none"
	FloatLeft            = "left"
	FloatRight           = "right"
)

type Length string

const (
	LengthAuto Length = "auto"
)

func LengthOf(val uint16, unit Unit) Length {
	return Length(strconv.FormatUint(uint64(val), 10)) + Length(unit)
}

type ResizeKind string

const (
	ResizeNone       ResizeKind = "none"
	ResizeBoth                  = "both"
	ResizeHorizontal            = "horizontal"
	ResizeVertical              = "vertical"
)

type Side string

const (
	Bottom Side = "Bottom"
	Left        = "Left"
	Right       = "Right"
	Top         = "Top"
)

var SideAll []Side = []Side{Bottom, Left, Right, Top}

type Unit string

const (
	// absolute units
	CM = "cm"
	MM = "mm"
	IN = "in" // 1in = 96px = 2.54cm
	PX = "px" // 1px = 1/96in
	PT = "pt" // 1pt = 1/72in
	PC = "pc" // 1pc = 1/6in
	//relative units
	EM   = "em"   // current font-size
	EX   = "ex"   // x-height of current font
	CH   = "ch"   // width of "0"
	REM  = "rem"  // root element font-size
	VW   = "vm"   // 1% of window width
	VH   = "vh"   // 1% of window height
	VMIN = "vmin" // min(vm, vh)
	VMAX = "vmax" // max(vm, vh)
	PCT  = "%"    // percentage of parent element
)

func fmtLength(val uint16, unit Unit) string {
	return strconv.FormatUint(uint64(val), 10) + string(unit)
}

func fmtUint16(val uint16) string {
	return strconv.FormatUint(uint64(val), 10)
}

type Field struct {
	Name  string
	Value string
}

type Attribute []Field

func side(name, val string, sides ...Side) Attribute {
	o := make(Attribute, len(sides))
	for i, side := range sides {
		o[i] = Field{
			Name:  name + string(side),
			Value: val,
		}
	}
	return o
}

func sideLong(base, name, val string, sides ...Side) Attribute {
	o := make(Attribute, len(sides))
	for i, side := range sides {
		o[i] = Field{
			Name:  base + string(side) + name,
			Value: val,
		}
	}
	return o
}

func BackgroundColor(color Color) Attribute {
	return Attribute{Field{
		Name:  "backgroundColor",
		Value: string(color),
	}}
}

func Border(width uint16, unit Unit, style BorderStyleKind, color Color, sides ...Side) Attribute {
	val := fmtLength(width, unit) + " " + string(style) + " " + string(color)
	return side("border", val, sides...)
}

func BorderCollapse(val bool) Attribute {
	var str string
	if val {
		str = "collapse"
	} else {
		str = "separate"
	}
	return Attribute{Field{
		Name:  "borderCollapse",
		Value: str,
	}}
}

func BorderColor(color Color, sides ...Side) Attribute {
	return sideLong("border", "Color", string(color), sides...)
}

func BorderRadius(val uint16, unit Unit, corners ...Corner) Attribute {
	o := make(Attribute, len(corners))
	for i, corner := range corners {
		o[i] = Field{
			Name:  "border" + string(corner) + "Radius",
			Value: fmtLength(val, unit),
		}
	}
	return o
}

func BorderStyle(style BorderStyleKind, sides ...Side) Attribute {
	return sideLong("border", "Style", string(style), sides...)
}

func BorderWidth(width uint16, unit Unit, sides ...Side) Attribute {
	return sideLong("border", "Width", fmtLength(width, unit), sides...)
}

func Cursor(val CursorKind) Attribute {
	return Attribute{Field{
		Name:  "cursor",
		Value: string(val),
	}}
}

func Display(val DisplayKind) Attribute {
	return Attribute{Field{
		Name:  "display",
		Value: string(val),
	}}
}

func Float(val FloatKind) Attribute {
	return Attribute{Field{
		Name:  "cssFloat",
		Value: string(val),
	}}
}

// "max-content" seems like a more sensible default, rather than the ambiguous "auto".
func gridVal(val Length) string {
	if val == LengthAuto {
		return "max-content"
	}
	return string(val)
}

func Grid(rows []Length, cols []Length) Attribute {
	return append(
		GridCols(cols...),
		GridRows(rows...)...,
	)
}

func GridCols(cols ...Length) Attribute {
	str := gridVal(cols[0])
	for i, n := 1, len(cols); i < n; i++ {
		str += " " + gridVal(cols[i])
	}
	return Attribute{Field{
		Name:  "gridTemplateColumns",
		Value: str,
	}}
}

func GridRows(rows ...Length) Attribute {
	str := gridVal(rows[0])
	for i, n := 1, len(rows); i < n; i++ {
		str += " " + gridVal(rows[i])
	}
	return Attribute{Field{
		Name:  "gridTemplateRows",
		Value: str,
	}}
}

// Indices start at 0.
func GridArea(rowStart, rowSpan, colStart, colSpan uint16) Attribute {
	return append(
		GridAreaRow(rowStart, rowSpan),
		GridAreaCol(colStart, colSpan)...,
	)
}

func GridAreaCol(start, span uint16) Attribute {
	return Attribute{Field{
		Name:  "gridColumnStart",
		Value: fmtUint16(start + 1),
	}, Field{
		Name:  "gridColumnEnd",
		Value: "span " + fmtUint16(span),
	}}
}

func GridAreaRow(start, span uint16) Attribute {
	return Attribute{Field{
		Name:  "gridRowStart",
		Value: fmtUint16(start + 1),
	}, Field{
		Name:  "gridRowEnd",
		Value: "span " + fmtUint16(span),
	}}
}

func Height(val uint16, unit Unit) Attribute {
	return Attribute{Field{
		Name:  "height",
		Value: fmtLength(val, unit),
	}}
}

func HeightMax(val uint16, unit Unit) Attribute {
	return Attribute{Field{
		Name:  "maxHeight",
		Value: fmtLength(val, unit),
	}}
}

func HeightMin(val uint16, unit Unit) Attribute {
	return Attribute{Field{
		Name:  "minHeight",
		Value: fmtLength(val, unit),
	}}
}

func Margin(val uint16, unit Unit, sides ...Side) Attribute {
	return side("margin", fmtLength(val, unit), sides...)
}

func Padding(val uint16, unit Unit, sides ...Side) Attribute {
	return side("padding", fmtLength(val, unit), sides...)
}

func Resize(val ResizeKind) Attribute {
	return Attribute{Field{
		Name:  "resize",
		Value: string(val),
	}}
}

func TextAlign(val Align) Attribute {
	return Attribute{Field{
		Name:  "textAlign",
		Value: string(val),
	}}
}

func Width(val uint16, unit Unit) Attribute {
	return Attribute{Field{
		Name:  "width",
		Value: fmtLength(val, unit),
	}}
}

func WidthMax(val uint16, unit Unit) Attribute {
	return Attribute{Field{
		Name:  "maxWidth",
		Value: fmtLength(val, unit),
	}}
}

func WidthMin(val uint16, unit Unit) Attribute {
	return Attribute{Field{
		Name:  "minWidth",
		Value: fmtLength(val, unit),
	}}
}
