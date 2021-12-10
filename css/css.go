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
	DarkSlateGrey        = "darkslategrey"
	Salmon               = "salmon"
	Thistle              = "thistle"
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
	DisplayNone                    = "none"
	DisplayTable                   = "table"
	DisplayTableCell               = "table-cell"
	DisplayTableRow                = "table-row"
)

type FloatKind string

const (
	FloatNone  FloatKind = "none"
	FloatLeft            = "left"
	FloatRight           = "right"
)

type FontStyleKind string

const (
	FontStyleItalic FontStyleKind = "italic"
	FontStyleNormal               = "normal"
)

type FontWeightKind string

const (
	FontWeightBold   FontWeightKind = "bold"
	FontWeightNormal                = "normal"
)

type Length string

const (
	LengthAuto Length = "auto"
)

func LengthOf(val uint16, unit Unit) Length {
	return Length(strconv.FormatUint(uint64(val), 10)) + Length(unit)
}

type PositionKind string

const (
	PositionAbsolute PositionKind = "absolute"
	PositionFixed                 = "fixed"
	PositionRelative              = "relative"
	PositionStatic                = "static"
	PositionSticky                = "sticky"
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

type SpaceKind string

const (
	SpaceNormal  SpaceKind = "normal"
	SpaceNoWrap            = "nowrap"
	SpacePre               = "pre"
	SpacePreLine           = "pre-line"
	SpacePreWrap           = "pre-wrap"
)

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

type VAlign string

const (
	VBaseline   VAlign = "baseline"
	VBottom            = "bottom"
	VMiddle            = "middle"
	VTop               = "top"
	VTextBottom        = "text-bottom"
	VTextTop           = "text-top"
	VSub               = "sub"
	VSuper             = "super"
)

func fmtLength(val uint16, unit Unit) string {
	return strconv.FormatUint(uint64(val), 10) + string(unit)
}

func fmtUint16(val uint16) string {
	return strconv.FormatUint(uint64(val), 10)
}

type Style map[string]string

// MakeStyle returns a new style that unites all argument styles.
func MakeStyle(src ...Style) Style {
	x := make(Style)
	for _, s := range src {
		for k, v := range s {
			x[k] = v
		}
	}
	return x
}

// Fork is a shorthand for MakeStyle(x, src...)
func (x Style) Fork(src ...Style) Style {
	return MakeStyle(append(src, x)...)
}

// Set includes the argument styles into the target.
func (x Style) Set(src ...Style) {
	for _, s := range src {
		for k, v := range s {
			x[k] = v
		}
	}
}

func side(name, val string, sides ...Side) Style {
	o := make(Style, len(sides))
	for _, side := range sides {
		k := name + string(side)
		o[k] = val
	}
	return o
}

func sideLong(base, name, val string, sides ...Side) Style {
	o := make(Style, len(sides))
	for _, side := range sides {
		k := base + string(side) + name
		o[k] = val
	}
	return o
}

func AlignText(val Align) Style {
	return Style{"textAlign": string(val)}
}

func AlignVertical(val VAlign) Style {
	return Style{"verticalAlign": string(val)}
}

func AlignVerticalL(val uint16, unit Unit) Style {
	return Style{"verticalAlign": fmtLength(val, unit)}
}

func BackgroundColor(color Color) Style {
	return Style{"backgroundColor": string(color)}
}

func Border(width uint16, unit Unit, style BorderStyleKind, color Color, sides ...Side) Style {
	val := fmtLength(width, unit) + " " + string(style) + " " + string(color)
	return side("border", val, sides...)
}

func BorderCollapse(val bool) Style {
	var str string
	if val {
		str = "collapse"
	} else {
		str = "separate"
	}
	return Style{"borderCollapse": str}
}

func BorderColor(color Color, sides ...Side) Style {
	return sideLong("border", "Color", string(color), sides...)
}

func BorderRadius(val uint16, unit Unit, corners ...Corner) Style {
	o := make(Style, len(corners))
	for _, corner := range corners {
		k := "border" + string(corner) + "Radius"
		o[k] = fmtLength(val, unit)
	}
	return o
}

func BorderStyle(style BorderStyleKind, sides ...Side) Style {
	return sideLong("border", "Style", string(style), sides...)
}

func BorderWidth(width uint16, unit Unit, sides ...Side) Style {
	return sideLong("border", "Width", fmtLength(width, unit), sides...)
}

func Cursor(val CursorKind) Style {
	return Style{"cursor": string(val)}
}

func Display(val DisplayKind) Style {
	return Style{"display": string(val)}
}

func Float(val FloatKind) Style {
	return Style{"cssFloat": string(val)}
}

// Font sets the css font family. Quotes the name.
func Font(name string) Style {
	return Style{"fontFamily": "\"" + name + "\""}
}

func FontSize(size uint16, unit Unit) Style {
	return Style{"fontSize": fmtLength(size, unit)}
}

func FontStyle(val FontStyleKind) Style {
	return Style{"fontStyle": string(val)}
}

func FontWeight(val FontWeightKind) Style {
	return Style{"fontWeight": string(val)}
}

// "max-content" seems like a more sensible default, rather than the ambiguous "auto".
func gridVal(val Length) string {
	if val == LengthAuto {
		return "max-content"
	}
	return string(val)
}

func Grid(rows []Length, cols []Length) Style {
	return MakeStyle(
		GridCols(cols...),
		GridRows(rows...),
	)
}

func GridCols(cols ...Length) Style {
	str := gridVal(cols[0])
	for i, n := 1, len(cols); i < n; i++ {
		str += " " + gridVal(cols[i])
	}
	return Style{"gridTemplateColumns": str}
}

func GridRows(rows ...Length) Style {
	str := gridVal(rows[0])
	for i, n := 1, len(rows); i < n; i++ {
		str += " " + gridVal(rows[i])
	}
	return Style{"gridTemplateRows": str}
}

// Indices start at 0.
func GridArea(rowStart, rowSpan, colStart, colSpan uint16) Style {
	return MakeStyle(
		GridAreaRow(rowStart, rowSpan),
		GridAreaCol(colStart, colSpan),
	)
}

func GridAreaCol(start, span uint16) Style {
	return Style{
		"gridColumnStart": fmtUint16(start + 1),
		"gridColumnEnd":   "span " + fmtUint16(span),
	}
}

func GridAreaRow(start, span uint16) Style {
	return Style{
		"gridRowStart": fmtUint16(start + 1),
		"gridRowEnd":   "span " + fmtUint16(span),
	}
}

func Height(val uint16, unit Unit) Style {
	return Style{"height": fmtLength(val, unit)}
}

func HeightMax(val uint16, unit Unit) Style {
	return Style{"maxHeight": fmtLength(val, unit)}
}

func HeightMin(val uint16, unit Unit) Style {
	return Style{"minHeight": fmtLength(val, unit)}
}

func Margin(val uint16, unit Unit, sides ...Side) Style {
	return side("margin", fmtLength(val, unit), sides...)
}

func OutlineStyle(val BorderStyleKind) Style {
	return Style{"outlineStyle": string(val)}
}

func Padding(val uint16, unit Unit, sides ...Side) Style {
	return side("padding", fmtLength(val, unit), sides...)
}

func Position(val PositionKind) Style {
	return Style{"position": string(val)}
}

func Resize(val ResizeKind) Style {
	return Style{"resize": string(val)}
}

func TabSize(val uint8) Style {
	return Style{"tabSize": strconv.FormatUint(uint64(val), 10)}
}

func TextColor(val Color) Style {
	return Style{"color": string(val)}
}

func TextLineHeight(coef float64) Style {
	return Style{"lineHeight": strconv.FormatFloat(coef, 'f', 1, 64)}
}

func Translate(x int16, unitX Unit, y int16, unitY Unit) Style {
	valX := strconv.Itoa(int(x)) + string(unitX)
	valY := strconv.Itoa(int(y)) + string(unitY)
	return Style{"transform": "translate(" + valX + "," + valY + ")"}
}

func WhiteSpace(val SpaceKind) Style {
	return Style{"whiteSpace": string(val)}
}

func Width(val uint16, unit Unit) Style {
	return Style{"width": fmtLength(val, unit)}
}

func WidthMax(val uint16, unit Unit) Style {
	return Style{"maxWidth": fmtLength(val, unit)}
}

func WidthMin(val uint16, unit Unit) Style {
	return Style{"minWidth": fmtLength(val, unit)}
}

// X sets the position on the horizontal axis
func X(val uint16, unit Unit) Style {
	return Style{"left": fmtLength(val, unit)}
}

// Y sets the position on the vertical axis
func Y(val uint16, unit Unit) Style {
	return Style{"top": fmtLength(val, unit)}
}
