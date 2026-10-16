// Package css provides CSS Go definitions.
package css

import (
	"strconv"

	"github.com/blitz-frost/wasm"
)

type Align string

const (
	AlignLeft    Align = "left"
	AlignRight   Align = "right"
	AlignCenter  Align = "center"
	AlignJustify Align = "justify"
)

type BorderStyleKind string

const (
	BorderNone   BorderStyleKind = "none"
	BorderHidden BorderStyleKind = "hidden"
	BorderDotted BorderStyleKind = "dotted"
	BorderDashed BorderStyleKind = "dashed"
	BorderSolid  BorderStyleKind = "solid"
	BorderDouble BorderStyleKind = "double"
	BorderGroove BorderStyleKind = "groove"
	BorderRidge  BorderStyleKind = "ridge"
	BorderInset  BorderStyleKind = "inset"
	BorderOutset BorderStyleKind = "outset"
)

type Color string

const (
	Transparent    Color = "transparent"
	Black          Color = "black"
	CadetBlue      Color = "cadetblue"
	CornflowerBlue Color = "cornflowerblue"
	Cyan           Color = "cyan"
	GhostWhite     Color = "ghostwhite"
	Grey           Color = "grey"
	LightGrey      Color = "lightgrey"
	DarkSlateGrey  Color = "darkslategrey"
	LightSlateGrey Color = "lightslategrey"
	Salmon         Color = "salmon"
	Thistle        Color = "thistle"
	White          Color = "white"
	WhiteSmoke     Color = "whitesmoke"
)

type Corner string

const (
	BottomLeft  Corner = "BottomLeft"
	BottomRight Corner = "BottomRight"
	TopLeft     Corner = "TopLeft"
	TopRight    Corner = "TopRight"
)

var CornerAll []Corner = []Corner{BottomLeft, BottomRight, TopLeft, TopRight}

type CursorKind string

const (
	CursorAlias        CursorKind = "alias"
	CursorScroll       CursorKind = "all-scroll"
	CursorCell         CursorKind = "cell"
	CursorContext      CursorKind = "context-menu"
	CursorCopy         CursorKind = "copy"
	CursorCross        CursorKind = "crosshair"
	CursorDefault      CursorKind = "default"
	CursorHelp         CursorKind = "help"
	CursorMove         CursorKind = "move"
	CursorNodrop       CursorKind = "no-drop"
	CursorNone         CursorKind = "none"
	CursorInvalid      CursorKind = "not-allowed"
	CursorPointer      CursorKind = "pointer"
	CursorProgress     CursorKind = "progress"
	CursorResizeCol    CursorKind = "col-resize"
	CursorResizeRow    CursorKind = "row-resize"
	CursorResizeE      CursorKind = "e-resize"
	CursorResizeEW     CursorKind = "ew-resize"
	CursorResizeN      CursorKind = "n-resize"
	CursorResizeNE     CursorKind = "ne-resize"
	CursorResizeNESW   CursorKind = "nesw-resize"
	CursorResizeNS     CursorKind = "ns-resize"
	CursorResizeNW     CursorKind = "nw-resize"
	CursorResizeNWSE   CursorKind = "nwse-resize"
	CursorResizeS      CursorKind = "s-resize"
	CursorResizeSE     CursorKind = "se-resize"
	CursorResizeSW     CursorKind = "sw-resize"
	CursorResizeW      CursorKind = "w-resize"
	CursorText         CursorKind = "text"
	CursorTextVertical CursorKind = "vertical-text"
	CursorWait         CursorKind = "wait"
	CursorZoomIn       CursorKind = "zoom-in"
	CursorZoomOut      CursorKind = "zoom-out"
)

type DisplayKind string

const (
	DisplayBlock       DisplayKind = "block"
	DisplayFlex        DisplayKind = "flex"
	DisplayGrid        DisplayKind = "grid"
	DisplayInline      DisplayKind = "inline"
	DisplayInlineBlock DisplayKind = "inline-block"
	DisplayInlineFlex  DisplayKind = "inline-flex"
	DisplayInlineGrid  DisplayKind = "inline-grid"
	DisplayInlineTable DisplayKind = "inline-table"
	DisplayNone        DisplayKind = "none"
	DisplayTable       DisplayKind = "table"
	DisplayTableCell   DisplayKind = "table-cell"
	DisplayTableRow    DisplayKind = "table-row"
)

type FloatKind string

const (
	FloatNone  FloatKind = "none"
	FloatLeft  FloatKind = "left"
	FloatRight FloatKind = "right"
)

type FontStyleKind string

const (
	FontStyleItalic FontStyleKind = "italic"
	FontStyleNormal FontStyleKind = "normal"
)

type FontWeightKind string

const (
	FontWeightBold   FontWeightKind = "bold"
	FontWeightNormal FontWeightKind = "normal"
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
	PositionFixed    PositionKind = "fixed"
	PositionRelative PositionKind = "relative"
	PositionStatic   PositionKind = "static"
	PositionSticky   PositionKind = "sticky"
)

type ResizeKind string

const (
	ResizeNone       ResizeKind = "none"
	ResizeBoth       ResizeKind = "both"
	ResizeHorizontal ResizeKind = "horizontal"
	ResizeVertical   ResizeKind = "vertical"
)

type Side string

const (
	Bottom Side = "Bottom"
	Left   Side = "Left"
	Right  Side = "Right"
	Top    Side = "Top"
)

var SideAll []Side = []Side{Bottom, Left, Right, Top}

type SpaceKind string

const (
	SpaceNormal  SpaceKind = "normal"
	SpaceNoWrap  SpaceKind = "nowrap"
	SpacePre     SpaceKind = "pre"
	SpacePreLine SpaceKind = "pre-line"
	SpacePreWrap SpaceKind = "pre-wrap"
)

type Unit string

const (
	// absolute units
	CM Unit = "cm"
	MM Unit = "mm"
	IN Unit = "in" // 1in = 96px = 2.54cm
	PX Unit = "px" // 1px = 1/96in
	PT Unit = "pt" // 1pt = 1/72in
	PC Unit = "pc" // 1pc = 1/6in
	//relative units
	EM   Unit = "em"   // current font-size
	EX   Unit = "ex"   // x-height of current font
	CH   Unit = "ch"   // width of "0"
	REM  Unit = "rem"  // root element font-size
	VW   Unit = "vm"   // 1% of window width
	VH   Unit = "vh"   // 1% of window height
	VMIN Unit = "vmin" // min(vm, vh)
	VMAX Unit = "vmax" // max(vm, vh)
	PCT  Unit = "%"    // percentage of parent element
)

type VAlign string

const (
	VBaseline   VAlign = "baseline"
	VBottom     VAlign = "bottom"
	VMiddle     VAlign = "middle"
	VTop        VAlign = "top"
	VTextBottom VAlign = "text-bottom"
	VTextTop    VAlign = "text-top"
	VSub        VAlign = "sub"
	VSuper      VAlign = "super"
)

func fmtLength(val uint16, unit Unit) string {
	return strconv.FormatUint(uint64(val), 10) + string(unit)
}

func fmtUint16(val uint16) string {
	return strconv.FormatUint(uint64(val), 10)
}

type Style struct {
	wasm.Object
}

func (x Style) AlignTextSet(val Align) {
	x.Set("textAlign", string(val))
}

func (x Style) AlignVerticalSet(val VAlign) {
	x.Set("verticalAlign", string(val))
}

func (x Style) AlignVerticalLSet(val uint16, unit Unit) {
	x.Set("verticalAlign", fmtLength(val, unit))
}

func (x Style) BackgroundColorSet(color Color) {
	x.Set("backgroundColor", string(color))
}

func (x Style) BorderSet(width uint16, unit Unit, style BorderStyleKind, color Color, sides ...Side) {
	val := fmtLength(width, unit) + " " + string(style) + " " + string(color)
	x._sidesSet("border", val, sides...)
}

func (x Style) BorderCollapseSet(val bool) {
	var str string
	if val {
		str = "collapse"
	} else {
		str = "separate"
	}
	x.Set("borderCollapse", str)
}

func (x Style) BorderColorSet(color Color, sides ...Side) {
	x._sidesLongSet("border", "Color", string(color), sides...)
}

func (x Style) BorderRadiusSet(val uint16, unit Unit, corners ...Corner) {
	for _, corner := range corners {
		prop := "border" + string(corner) + "Radius"
		value := fmtLength(val, unit)
		x.Set(prop, value)
	}
}

func (x Style) BorderStyleSet(style BorderStyleKind, sides ...Side) {
	x._sidesLongSet("border", "Style", string(style), sides...)
}

func (x Style) BorderWidthSet(width uint16, unit Unit, sides ...Side) {
	x._sidesLongSet("border", "Width", fmtLength(width, unit), sides...)
}

func (x Style) CursorSet(val CursorKind) {
	x.Set("cursor", string(val))
}

func (x Style) DisplaySet(val DisplayKind) {
	x.Set("display", string(val))
}

func (x Style) FloatSet(val FloatKind) {
	x.Set("cssFloat", string(val))
}

// Font sets the css font family. Quotes the name.
func (x Style) FontSet(name string) {
	x.Set("fontFamily", "\""+name+"\"")
}

func (x Style) FontSizeSet(size uint16, unit Unit) {
	x.Set("fontSize", fmtLength(size, unit))
}

func (x Style) FontStyleSet(val FontStyleKind) {
	x.Set("fontStyle", string(val))
}

func (x Style) FontWeightSet(val FontWeightKind) {
	x.Set("fontWeight", string(val))
}

func (x Style) GridSet(rows, cols []Length) {
	x.GridRowsSet(rows...)
	x.GridColsSet(cols...)
}

func (x Style) GridColsSet(cols ...Length) {
	str := gridVal(cols[0])
	for _, col := range cols[1:] {
		str += " " + gridVal(col)
	}
	x.Set("gridTemplateColumns", str)
}

func (x Style) GridRowsSet(rows ...Length) {
	str := gridVal(rows[0])
	for _, row := range rows[1:] {
		str += " " + gridVal(row)
	}
	x.Set("gridTemplateRows", str)
}

// Indices start at 0.
func (x Style) GridAreaSet(rowStart, rowSpan, colStart, colSpan uint16) {
	x.GridAreaRowSet(rowStart, rowSpan)
	x.GridAreaColSet(colStart, colSpan)
}

func (x Style) GridAreaColSet(start, span uint16) {
	x.Set("gridColumnStart", fmtUint16(start+1))
	x.Set("gridColumnEnd", "span "+fmtUint16(span))
}

func (x Style) GridAreaRowSet(start, span uint16) {
	x.Set("gridRowStart", fmtUint16(start+1))
	x.Set("gridRowEnd", "span "+fmtUint16(span))

}

func (x Style) HeightSet(val uint16, unit Unit) {
	x.Set("height", fmtLength(val, unit))
}

func (x Style) HeightMaxSet(val uint16, unit Unit) {
	x.Set("maxHeight", fmtLength(val, unit))
}

func (x Style) HeightMinSet(val uint16, unit Unit) {
	x.Set("minHeight", fmtLength(val, unit))
}

func (x Style) MarginSet(val uint16, unit Unit, sides ...Side) {
	x._sidesSet("margin", fmtLength(val, unit), sides...)
}

func (x Style) OutlineStyleSet(val BorderStyleKind) {
	x.Set("outlineStyle", string(val))
}

func (x Style) PaddingSet(val uint16, unit Unit, sides ...Side) {
	x._sidesSet("padding", fmtLength(val, unit), sides...)
}

func (x Style) PositionSet(val PositionKind) {
	x.Set("position", string(val))
}

func (x Style) ResizeSet(val ResizeKind) {
	x.Set("resize", string(val))
}

func (x Style) TabSizeSet(val uint8) {
	x.Set("tabSize", strconv.FormatUint(uint64(val), 10))
}

func (x Style) TextColorSet(val Color) {
	x.Set("color", string(val))
}

func (x Style) TextLineHeightSet(coef float64) {
	x.Set("lineHeight", strconv.FormatFloat(coef, 'f', 1, 64))
}

func (x Style) TranslateSet(X int16, unitX Unit, Y int16, unitY Unit) {
	valX := strconv.Itoa(int(X)) + string(unitX)
	valY := strconv.Itoa(int(Y)) + string(unitY)
	x.Set("transform", "translate("+valX+","+valY+")")
}

func (x Style) WhiteSpaceSet(val SpaceKind) {
	x.Set("whiteSpace", string(val))
}

func (x Style) WidthSet(val uint16, unit Unit) {
	x.Set("width", fmtLength(val, unit))
}

func (x Style) WidthMaxSet(val uint16, unit Unit) {
	x.Set("maxWidth", fmtLength(val, unit))
}

func (x Style) WidthMinSet(val uint16, unit Unit) {
	x.Set("minWidth", fmtLength(val, unit))
}

// X is the position on the horizontal axis, originating on the left.
func (x Style) XSet(val uint16, unit Unit) {
	x.Set("left", fmtLength(val, unit))
}

// Y is the position on the vertical axis, originating at the top.
func (x Style) YSet(val uint16, unit Unit) {
	x.Set("top", fmtLength(val, unit))
}

func (x Style) _sidesSet(name, val string, sides ...Side) {
	for _, side := range sides {
		prop := name + string(side)
		x.Set(prop, val)
	}
}

func (x Style) _sidesLongSet(base, name, val string, sides ...Side) {
	for _, side := range sides {
		prop := base + string(side) + name
		x.Set(prop, val)
	}
}

// "max-content" seems like a more sensible default, rather than the ambiguous "auto".
func gridVal(val Length) string {
	if val == LengthAuto {
		return "max-content"
	}
	return string(val)
}
