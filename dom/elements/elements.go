// Package elements provides definitions for common DOM elements.
package elements

import (
	"syscall/js"

	"github.com/blitz-frost/wasm/dom"
	"github.com/blitz-frost/wasm/media"
)

var global = js.Global()
var doc = global.Get("document")

type Element = dom.Element

type Button struct {
	Element
}

func MakeButton() Button {
	return Button{Element{doc.Call("createElement", "button")}}
}

// A Cell wraps a DOM td
type Cell struct {
	Element
}

func MakeCell() Cell {
	return Cell{Element{doc.Call("createElement", "td")}}
}

func (x Cell) Index() int {
	return x.Get("cellIndex").Int()
}

func (x Cell) SpanColSet(n int) {
	x.Set("colSpan", n)
}

func (x Cell) SpanRowSet(n int) {
	x.Set("rowSpan", n)
}

func (x Cell) Row() Row {
	return Row{Element{x.Get("parentElement")}}
}

type Checkbox struct {
	Element
}

func MakeCheckbox() Checkbox {
	e := Element{doc.Call("createElement", "input")}
	e.Call("setAttribute", "type", "checkbox")
	return Checkbox{e}
}

func (x Checkbox) Checked() bool {
	return x.Get("checked").Bool()
}

func (x Checkbox) CheckedSet(v bool) {
	x.Set("checked", v)
}

func (x Checkbox) Default() bool {
	return x.Get("defaultChecked").Bool()
}

func (x Checkbox) DefaultSet(v bool) {
	x.Set("defaultChecked", v)
}

func (x Checkbox) Explicit() bool {
	return !x.Get("indeterminate").Bool()
}

func (x Checkbox) ExplicitSet(v bool) {
	x.Set("indeterminate", !v)
}

type Div struct {
	Element
}

func MakeDiv() Div {
	return Div{Element{doc.Call("createElement", "div")}}
}

type Image struct {
	Element
}

func MakeImage() Image {
	return Image{Element{doc.Call("createElement", "img")}}
}

func (x Image) Src() string {
	return x.Call("getAttribute", "src").String()
}

func (x Image) SrcSet(s string) {
	x.Call("setAttribute", "src", s)
}

type Option struct {
	Element
}

func MakeOption(val string) Option {
	x := Option{Element{doc.Call("createElement", "option")}}
	x.ValueSet(val)
	return x
}

func (x Option) Value() string {
	return x.Element.Get("value").String()
}

func (x Option) ValueSet(val string) {
	x.Set("value", val)
	x.Set("text", val)
}

type Para struct {
	Element
}

func MakePara() Para {
	return Para{Element{doc.Call("createElement", "p")}}
}

// A Row wraps a DOM tr
type Row struct {
	Element
}

func MakeRow() Row {
	return Row{Element{doc.Call("createElement", "tr")}}
}

func (x Row) Add(pos int, cell ...Cell) {
	jsCell := x.Get("cells").Index(pos)
	for _, c := range cell {
		x.Call("insertBefore", c.Element.Value, jsCell)
	}
}

func (x Row) Append(cell ...Cell) {
	for _, c := range cell {
		x.Call("appendChild", c.Element.Value)
	}
}

// Cell returns the row's i-th cell, starting at 0.
func (x Row) Cell(i int) Cell {
	return Cell{Element{x.Get("cells").Index(i)}}
}

func (x Row) Delete(i int) {
	x.Call("deleteCell", i)
}

// Index returns the row's position in the table that contains it.
func (x Row) Index() int {
	return x.Get("rowIndex").Int()
}

// Len returns the row's number of cells.
func (x Row) Len() int {
	return x.Get("cells").Length()
}

func (x Row) Table() Table {
	return Table{Element{x.Get("parentElement")}}
}

type Select struct {
	Element
}

func MakeSelect(opt ...Option) Select {
	x := Select{Element{doc.Call("createElement", "select")}}
	x.Append(opt...)
	return x
}

func (x Select) Add(pos int, opt ...Option) {
	for i, op := range opt {
		x.Call("add", op.Element.Value, pos+i)
	}
}

func (x Select) Append(opt ...Option) {
	for _, op := range opt {
		x.Call("add", op.Element.Value)
	}
}

// IndexSet sets the currently active option.
func (x Select) IndexSet(i int) {
	x.Set("selectedIndex", i)
}

// Get returns the value of the currently selected option.
func (x Select) Value() string {
	return x.Get("value").String()
}

// Set attempts to set the active selected option based on the given value.
// If no option has that value, the active option will be empty.
func (x Select) ValueSet(val string) {
	x.Set("value", val)
}

func (x Select) Len() int {
	return x.Element.Get("options").Length()
}

type Table struct {
	Element
}

func MakeTable() Table {
	return Table{Element{doc.Call("createElement", "table")}}
}

func (x Table) Add(i int, row ...Row) {
	jsRow := x.Get("rows").Index(i)
	for _, r := range row {
		x.Call("insertBefore", r.Element.Value, jsRow)
	}
}

func (x Table) Append(row ...Row) {
	for _, r := range row {
		x.Call("appendChild", r.Element.Value)
	}
}

// Clear deletes all rows from the table.
func (x Table) Clear() {
	n := x.Len()
	for i := 0; i < n; i++ {
		x.Call("deleteRow", -1)
	}
}

// Delete removes the specified row from the table.
func (x Table) Delete(i int) {
	x.Call("deleteRow", i)
}

// Len returns the table's number of rows.
func (x Table) Len() int {
	return x.Get("rows").Length()
}

// Row returns the table's ith row, starting at 0.
func (x Table) Row(i int) Row {
	return Row{Element{x.Get("rows").Index(i)}}
}

type TextArea struct {
	Element
}

func MakeTextArea() TextArea {
	return TextArea{Element{doc.Call("createElement", "textarea")}}
}

func (x TextArea) PlaceholderSet(s string) {
	x.Set("placeholder", s)
}

func (x TextArea) RowsSet(n int) {
	x.Set("rows", n)
}

func (x TextArea) Text() string {
	return x.Get("value").String()
}

func (x TextArea) TextSet(s string) {
	x.Set("value", s)
}

type Video struct {
	Element
}

func MakeVideo() Video {
	return Video{Element{doc.Call("createElement", "video")}}
}

func (x Video) AutoPlay() bool {
	return x.Get("autoplay").Bool()
}

func (x Video) AutoPlaySet(v bool) {
	x.Set("autoplay", v)
}

func (x Video) Muted() bool {
	return x.Get("muted").Bool()
}

func (x Video) MutedSet(v bool) {
	x.Set("muted", v)
}

func (x Video) SourceStream() media.Stream {
	v := x.Get("srcObject")
	return media.AsStream(v)
}

func (x Video) SourceStreamSet(stream media.Stream) {
	v := stream.Js()
	x.Set("srcObject", v)
}

func (x Video) SourceUrl() string {
	return x.Get("src").String()
}

func (x Video) SourceUrlSet(url string) {
	x.Set("src", url)
}
