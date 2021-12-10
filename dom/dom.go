// Package dom wraps JS DOM elements.
package dom

import (
	"syscall/js"

	"github.com/blitz-frost/wasm/css"
)

var window js.Value = js.Global()
var doc js.Value = window.Get("document")
var Root Element = Element{doc.Call("getElementById", "root")}

type Base interface {
	Base() Element
}

// A Base represents a JS DOM element and forms the basis of this package.
// It wraps js.Value and gives access to all its funcionality.
type Element struct {
	js.Value // underlying JS object
}

// Add adds the given elements as subelements, at the given position.
func (x Element) Add(pos int, e ...Base) {
	jsVal := x.Get("children").Index(pos)
	for _, b := range e {
		x.Call("insertBefore", b.Base().Value, jsVal)
	}
}

// Append adds the given elements as final subelement.
func (x Element) Append(e ...Base) {
	for _, b := range e {
		x.Call("appendChild", b.Base().Value)
	}
}

func (x Element) Blur() {
	x.Call("blur")
}

func (x Element) Class() string {
	return x.Get("className").String()
}

func (x Element) ClassSet(name string) {
	x.Set("className", name)
}

// Delete removes the subelement at index i.
func (x Element) Delete(i int) {
	sub := x.Get("children").Index(i)
	sub.Call("remove")
}

func (x Element) EditableSet(t bool) {
	x.Set("contentEditable", t)
}

func (x Element) Focus() {
	x.Call("focus")
}

// Handle subscribes the given Handler to the specified event.
func (x Element) Handle(event EventName, h Handler) {
	x.Call("addEventListener", string(event), h.f)
}

// HandleRemove unsubscribes the given Handler from the specified event.
func (x Element) HandleRemove(event EventName, h Handler) {
	x.Call("removeEventListener", string(event), h.f)
}

func (x Element) Height() uint16 {
	return uint16(x.Get("offsetHeight").Float())
}

func (x Element) Id() string {
	return x.Get("id").String()
}

func (x Element) IdSet(id string) {
	x.Set("id", id)
}

// Len returns the number of subelement.
func (x Element) Len() int {
	return x.Get("children").Length()
}

// Next returns the next element in the same node.
// Returns an empty Element if there is none.
func (x Element) Next() Element {
	return Element{x.Get("nextElementSibling")}
}

// Previous returns the previous element in the same node.
// Returns an empty Element if there is none.
func (x Element) Previous() Element {
	return Element{x.Get("previousElementSibling")}
}

// Remove removes the specified subelements.
func (x Element) Remove(e ...Base) {
	for _, b := range e {
		x.Call("removeChild", b.Base().Value)
	}
}

// RemoveSelf removes the target Element from the dom
func (x Element) RemoveSelf() {
	x.Call("remove")
}

func (x Element) Replace(newElem, oldElem Base) {
	x.Call("replaceChild", newElem.Base().Value, oldElem.Base().Value)
}

func (x Element) SpellcheckSet(val bool) {
	x.Set("spellcheck", val)
}

// Style sets the value of the specified style component.
func (x Element) Style(style ...css.Style) {
	jsStyle := x.Get("style")
	for _, s := range style {
		for k, v := range s {
			jsStyle.Set(k, v)
		}
	}
}

func (x Element) Sub(i int) Element {
	return Element{(x.Get("children").Index(i))}
}

func (x Element) Super() Element {
	return Element{(x.Get("parentElement"))}
}

// Text returns the inner HTML text node value. Panics if x does not contain a text node.
func (x Element) Text() string {
	return x.Get("innerHTML").String()
}

// TextSet sets the inner HTML of x as a text node with the provided value.
func (x Element) TextSet(s string) {
	x.Set("innerHTML", s)
}

func (x Element) Width() uint16 {
	return uint16(x.Get("offsetWidth").Float())
}

func (x Element) Base() Element {
	return x
}

type Div struct {
	Element
}

func MakeDiv() Div {
	return Div{Element{doc.Call("createElement", "div")}}
}

type Para struct {
	Element
}

func MakePara() Para {
	return Para{Element{doc.Call("createElement", "p")}}
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

type Button struct {
	Element
}

func MakeButton() Button {
	return Button{Element{doc.Call("createElement", "button")}}
}

type Option struct {
	Element
}

func MakeOption() Option {
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

func MakeSelect() Select {
	return Select{Element{doc.Call("createElement", "select")}}
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
