// Package dom wraps JS DOM elements
package dom

import (
	"syscall/js"
)

var doc js.Value = js.Global().Get("document")
var Root Element = Element{doc.Call("getElementById", "root")}

// A Collection can pass itself as a set of DOM elements.
type Collection interface {
	Elements() []Element
}

// Compose is a helper function to implement the Element interface using multiple DOM objects.
func Compose(e ...Collection) []Element {
	if len(e) == 0 {
		return nil
	}

	r := e[0].Elements()
	for i := 1; i < len(e); i++ {
		r = append(r, e[i].Elements()...)
	}

	return r
}

// A Base represents a JS DOM element and forms the basis of this package.
// It wraps js.Value and gives access to all its funcionality.
type Element struct {
	js.Value // underlying JS object
}

func (x Element) Add(c Collection) {
	for _, v := range c.Elements() {
		x.Value.Call("appendChild", v.Value)
	}
}

func (x Element) Handle(event string, h Handler) {
	x.Set("on"+event, h.f)
}

func (x Element) Class() string {
	return x.Get("className").String()
}

func (x Element) ClassSet(name string) {
	x.Set("className", name)
}

func (x Element) EditableSet(t bool) {
	x.Set("contentEditable", t)
}

func (x Element) Focus() {
	x.Call("focus")
}

func (x Element) Id() string {
	return x.Get("id").String()
}

func (x Element) IdSet(id string) {
	x.Set("id", id)
}

// Style returns the current value of the specified style component.
func (x Element) Style(name string) string {
	return x.Get("style").Get(name).String()
}

// StyleSet sets the value of the specified style component.
func (x Element) StyleSet(name, value string) {
	x.Get("style").Set(name, value)
}

func (x Element) Sub(i int) Element {
	return Element{(x.Get("children").Index(i))}
}

func (x Element) SubLen() int {
	return x.Get("children").Length()
}

// Text returns the inner HTML text node value. Panics if x does not contain a text node.
func (x Element) Text() string {
	return x.Get("innerHTML").String()
}

// TextSet sets the inner HTML of x as a text node with the provided value.
func (x Element) TextSet(s string) {
	x.Set("innerHTML", s)
}

func (x Element) Elements() []Element {
	return []Element{x}
}

type Div struct {
	Element
}

func NewDiv() Div {
	return Div{Element{doc.Call("createElement", "div")}}
}

type TextArea struct {
	Element
}

func AsTextArea(e Element) TextArea {
	return TextArea{e}
}

func NewTextArea() TextArea {
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

type Button struct {
	Element
}

func NewButton() Button {
	return Button{Element{doc.Call("createElement", "button")}}
}

type Option struct {
	Element
}

func NewOption() Option {
	return Option{Element{doc.Call("createElement", "option")}}
}

func (x Option) Value() string {
	return x.Element.Get("value").String()
}

func (x Option) ValueSet(val string) {
	x.Set("value", val)
	x.Set("text", val)
}

type Select struct {
	Element
}

func NewSelect() Select {
	return Select{Element{doc.Call("createElement", "select")}}
}

func AsSelect(e Element) Select {
	return Select{e}
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

// Expand inserts a new option at the end of the list and returns it.
func (x Select) Expand() Option {
	r := NewOption()
	x.Call("add", r.Element.Value)
	return r
}

func (x Select) Len() int {
	return x.Element.Get("options").Length()
}

// A Cell wraps a DOM td
type Cell struct {
	Element
}

func AsCell(e Element) Cell {
	return Cell{e}
}

func (x Cell) SpanSet(n int) {
	x.Set("colSpan", n)
}

func (x Cell) Row() Row {
	return Row{Element{x.Get("parentElement")}}
}

// A Row wraps a DOM tr
type Row struct {
	Element
}

// Cell returns the row's i-th cell, starting at 0.
func (x Row) Cell(i int) Cell {
	return Cell{Element{x.Get("cells").Index(i)}}
}

// Expand inserts a new cell at the end of the row and returns it.
func (x Row) Expand() Cell {
	return Cell{Element{x.Call("insertCell", -1)}}
}

func (x Row) Index() int {
	return x.Get("rowIndex").Int()
}

// Len returns the row's number of cells.
func (x Row) Len() int {
	return x.Get("cells").Length()
}

type Table struct {
	Element
}

func AsTable(e Element) Table {
	return Table{e}
}

func NewTable() Table {
	return Table{Element{doc.Call("createElement", "table")}}
}

// Clear deletes all rows from the table.
func (x Table) Clear() {
	n := x.Len()
	for i := 0; i < n; i++ {
		x.Call("deleteRow", -1)
	}
}

// Expand inserts a new row at the bottom of the table and returns it.
func (x Table) Expand() Row {
	return Row{Element{x.Call("insertRow", -1)}}
}

// Len returns the table's number of rows.
func (x Table) Len() int {
	return x.Get("rows").Length()
}

// Row returns the table's ith row, starting at 0.
func (x Table) Row(i int) Row {
	return Row{Element{x.Get("rows").Index(i)}}
}
