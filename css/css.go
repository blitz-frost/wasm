// Package css provides CSS Go definitions.
package css

type Align string

const (
	AlignLeft    Align = "left"
	AlignRight         = "right"
	AlignCenter        = "center"
	AlignJustify       = "justify"
)

type Border string

const (
	BorderNone   Border = "none"
	BorderHidden        = "hidden"
	BorderDotted        = "dotted"
	BorderDashed        = "dashed"
	BorderSolid         = "solid"
	BorderDouble        = "double"
	BorderGroove        = "groove"
	BorderRidge         = "ridge"
	BorderInset         = "inset"
	BorderOutset        = "outset"
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

type Cursor string

const (
	CursorAlias        Cursor = "alias"
	CursorScroll              = "all-scroll"
	CursorCell                = "cell"
	CursorContext             = "context-menu"
	CursorCopy                = "copy"
	CursorCross               = "crosshair"
	CursorDefault             = "default"
	CursorHelp                = "help"
	CursorMove                = "move"
	CursorNodrop              = "no-drop"
	CursorNone                = "none"
	CursorInvalid             = "not-allowed"
	CursorPointer             = "pointer"
	CursorProgress            = "progress"
	CursorResizeCol           = "col-resize"
	CursorResizeRow           = "row-resize"
	CursorResizeE             = "e-resize"
	CursorResizeEW            = "ew-resize"
	CursorResizeN             = "n-resize"
	CursorResizeNE            = "ne-resize"
	CursorResizeNESW          = "nesw-resize"
	CursorResizeNS            = "ns-resize"
	CursorResizeNW            = "nw-resize"
	CursorResizeNWSE          = "nwse-resize"
	CursorResizeS             = "s-resize"
	CursorResizeSE            = "se-resize"
	CursorResizeSW            = "sw-resize"
	CursorResizeW             = "w-resize"
	CursorText                = "text"
	CursorTextVertical        = "vertical-text"
	CursorWait                = "wait"
	CursorZoomIn              = "zoom-in"
	CursorZoomOut             = "zoom-out"
)

type Float string

const (
	FloatNone  Float = "none"
	FloatLeft        = "left"
	FloatRight       = "right"
)

type Resize string

const (
	ResizeNone       Resize = "none"
	ResizeBoth              = "both"
	ResizeHorizontal        = "horizontal"
	ResizeVertical          = "vertical"
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
