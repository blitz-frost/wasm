// Package elements provides definitions for common DOM elements.
package elements

import (
	"github.com/blitz-frost/wasm"
	"github.com/blitz-frost/wasm/html"
	"github.com/blitz-frost/wasm/html/document"
	"github.com/blitz-frost/wasm/media"
)

type Button struct {
	html.Element
}

func ButtonMake() Button {
	v := document.ElementCreate("button")
	return Button{v}
}

// A Cell wraps a DOM td
type Cell struct {
	html.Element
}

func CellMake() Cell {
	v := document.ElementCreate("td")
	return Cell{v}
}

func (x Cell) Index() int {
	v := x.Get("cellIndex")
	return v.Int()
}

func (x Cell) InsertSelf(cell Cell) {
	x.Element.InsertSelf(cell.Element)
}

func (x Cell) SpanColSet(n int) {
	x.Set("colSpan", n)
}

func (x Cell) SpanRowSet(n int) {
	x.Set("rowSpan", n)
}

func (x Cell) Row() Row {
	v := x.Parent()
	return Row{v}
}

type Checkbox struct {
	html.Element
}

func CheckboxMake() Checkbox {
	v := document.ElementCreate("input")
	v.AttributeSet("type", "checkbox")
	return Checkbox{v}
}

func (x Checkbox) Checked() bool {
	v := x.Get("checked")
	return v.Bool()
}

func (x Checkbox) CheckedSet(v bool) {
	x.Set("checked", v)
}

func (x Checkbox) Default() bool {
	v := x.Get("defaultChecked")
	return v.Bool()
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
	html.Element
}

func DivMake() Div {
	v := document.ElementCreate("div")
	return Div{v}
}

type Image struct {
	html.Element
}

func ImageMake() Image {
	v := document.ElementCreate("img")
	return Image{v}
}

func (x Image) Src() (string, bool) {
	return x.Attribute("src")
}

func (x Image) SrcSet(s string) {
	x.AttributeSet("src", s)
}

type Option struct {
	html.Element
}

func OptionMake(val string) Option {
	v := document.ElementCreate("option")
	x := Option{v}
	x.ValueSet(val)
	return x
}

func (x Option) Value() string {
	return x.Get("value").String()
}

func (x Option) ValueSet(val string) {
	x.Set("value", val)
	x.Set("text", val)
}

type Para struct {
	html.Element
}

func ParaMake() Para {
	v := document.ElementCreate("p")
	return Para{v}
}

// A Row wraps a DOM tr
type Row struct {
	html.Element
}

func RowMake() Row {
	v := document.ElementCreate("tr")
	return Row{v}
}

func (x Row) Append(cell Cell) {
	x.Element.Append(cell.Element)
}

func (x Row) Cells() html.Collection {
	v := x.Get("cells")
	return html.Collection{wasm.Object(v)}
}

func (x Row) Delete(i int) {
	x.CallRaw("deleteCell", i)
}

// Index returns the row's position in the table that contains it.
func (x Row) Index() int {
	return x.Get("rowIndex").Int()
}

func (x Row) InsertSelf(row Row) {
	x.Element.InsertSelf(row.Element)
}

func (x Row) Table() Table {
	v := x.Parent()
	return Table{v}
}

type Select struct {
	html.Element
}

func SelectMake(opt ...Option) Select {
	v := document.ElementCreate("select")
	x := Select{v}
	for _, option := range opt {
		x.Append(option)
	}
	return x
}

// Activate sets the currently active option.
func (x Select) Activate(i int) {
	x.Set("selectedIndex", i)
}

func (x Select) Add(pos int, opt Option) {
	x.CallRaw("add", opt.Value(), pos)
}

func (x Select) Append(opt Option) {
	x.CallRaw("add", opt.Value())
}

func (x Select) Options() html.Collection {
	v := x.Get("options")
	return html.Collection{wasm.Object(v)}
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

type Table struct {
	html.Element
}

func TableMake() Table {
	v := document.ElementCreate("table")
	return Table{v}
}

func (x Table) Append(row Row) {
	x.Element.Append(row.Element)
}

// Delete removes the specified row from the table.
func (x Table) Delete(i int) {
	x.Call("deleteRow", i)
}

func (x Table) Rows() html.Collection {
	v := x.Get("rows")
	return html.Collection{wasm.Object(v)}
}

type TextArea struct {
	html.Element
}

func TextAreaMake() TextArea {
	v := document.ElementCreate("textarea")
	return TextArea{v}
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
	html.Element
}

func VideoMake() Video {
	v := document.ElementCreate("video")
	return Video{v}
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
	return media.Stream{wasm.Object(v)}
}

func (x Video) SourceStreamSet(stream media.Stream) {
	x.Set("srcObject", stream.V)
}

func (x Video) SourceUrl() string {
	return x.Get("src").String()
}

func (x Video) SourceUrlSet(url string) {
	x.Set("src", url)
}
