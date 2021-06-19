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

type FloatKind string

const (
	FloatNone  FloatKind = "none"
	FloatLeft            = "left"
	FloatRight           = "right"
)

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
	UnitCM Unit = "cm"
	UnitMM      = "mm"
	UnitIN      = "in" // 1in = 96px = 2.54cm
	UnitPX      = "px" // 1px = 1/96in
	UnitPT      = "pt" // 1pt = 1/72in
	UnitPC      = "pc" // 1pc = 1/6in
	//relative units
	UnitEM   = "em"   // current font-size
	UnitEX   = "ex"   // x-height of current font
	UnitCH   = "ch"   // width of "0"
	UnitREM  = "rem"  // root element font-size
	UnitVW   = "vm"   // 1% of window width
	UnitVH   = "vh"   // 1% of window height
	UnitVMIN = "vmin" // min(vm, vh)
	UnitVMAX = "vmax" // max(vm, vh)
	UnitPCT  = "%"    // percentage of parent element
)

func fmtLength(val uint16, unit Unit) string {
	return strconv.FormatUint(uint64(val), 10) + string(unit)
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

func Float(val FloatKind) Attribute {
	return Attribute{Field{
		Name:  "cssFloat",
		Value: string(val),
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
