// Package dom wraps JS DOM elements
package dom

import (
	"syscall/js"
)

var doc js.Value = js.Global().Get("document")
var Root Base = *newBase(doc.Call("getElementById", "root"))

// An Element can pass itself as a collection of DOM elements.
type Element interface {
	Elem() []*Base
}

// Compose is a helper function to implement the Element interface using multiple DOM objects.
func Compose(e ...Element) []*Base {
	if len(e) == 0 {
		return nil
	}

	r := e[0].Elem()
	for i := 1; i < len(e); i++ {
		r = append(r, e[i].Elem()...)
	}

	return r
}

// A Base represents a JS DOM element and forms the basis of this package.
// It wraps js.Value and gives access to all its funcionality.
type Base struct {
	js.Value // underlying JS object

	Id    StringAttr
	Class StringAttr
}

func newBase(obj js.Value) *Base {
	r := Base{Value: obj}

	r.Id = StringAttr{&r, "id"}
	r.Class = StringAttr{&r, "className"}

	return &r
}

func (x Base) Add(e Element) {
	for _, v := range e.Elem() {
		x.Value.Call("appendChild", v.Value)
	}
}

func (x Base) Handle(event string, h Handler) {
	x.Set(event, h.f)
}

// GetStyle returns the current value of the specified style component.
func (x Base) GetStyle(name string) string {
	return x.Get("style").Get(name).String()
}

// SetStyle sets the value of the specified style component.
func (x *Base) SetStyle(name, value string) {
	x.Get("style").Set(name, value)
}

// GetText returns the inner HTML text node value. Panics if x does not contain a text node.
func (x Base) GetText() string {
	return x.Get("innerHTML").String()
}

// SetText sets the inner HTML of x as a text node with the provided value.
func (x *Base) SetText(s string) {
	x.Set("innerHTML", s)
}

func (x *Base) Elem() []*Base {
	return []*Base{x}
}

type Div struct {
	*Base
}

func NewDiv() *Div {
	return &Div{newBase(doc.Call("createElement", "div"))}
}

type TextArea struct {
	*Base
}

func NewTextArea() *TextArea {
	return &TextArea{newBase(doc.Call("createElement", "textarea"))}
}

type Button struct {
	*Base
}

func NewButton() *Button {
	return &Button{newBase(doc.Call("createElement", "button"))}
}

type Option struct {
	*Base
}

func NewOption() *Option {
	return &Option{newBase(doc.Call("createElement", "option"))}
}

// Get returns the option's value.
func (x Option) Get() string {
	return x.Base.Get("value").String()
}

// Set sets the HTML option's value and text to v.
func (x *Option) Set(v string) {
	x.Base.Set("value", v)
	x.Base.Set("text", v)
}

type Select struct {
	*Base
	opts []Option
}

func NewSelect() *Select {
	return &Select{
		Base: newBase(doc.Call("createElement", "select")),
		opts: make([]Option, 0),
	}
}

// Set sets the currently active option.
func (x *Select) Set(index int) {
	x.Base.Set("selectedIndex", index)
}

// Get returns the value of the currently selected option.
func (x Select) Get() string {
	return x.Base.Get("value").String()
}

// Expand inserts a new option at the end of the list and returns it.
func (x *Select) Expand() *Option {
	n := len(x.opts)
	x.opts = append(x.opts, *NewOption())
	x.Call("add", x.opts[n].Base.Value)
	return &x.opts[n]
}

func (x *Select) Len() int {
	return len(x.opts)
}

// A Cell wraps a DOM td
type Cell struct {
	*Base
}

// A Row wraps a DOM tr
type Row struct {
	*Base
	cells []Cell // consistent with table
}

// Index returns the row's ith cell, starting at 0.
func (x Row) Index(i int) *Cell {
	return &x.cells[i]
}

// Len returns the row's number of cells.
func (x Row) Len() int {
	return len(x.cells)
}

// Expand inserts a new cell at the end of the row and returns it.
func (x *Row) Expand() *Cell {
	x.cells = append(x.cells, Cell{newBase(x.Call("insertCell", -1))})
	return &x.cells[len(x.cells)-1]
}

type Table struct {
	*Base
	rows []Row // more convenient to store a Go slice for types that hold more than a Base
}

func NewTable() *Table {
	return &Table{
		Base: newBase(doc.Call("createElement", "table")),
		rows: make([]Row, 0),
	}
}

// Index returns the table's ith row, starting at 0.
func (x Table) Index(i int) *Row {
	return &x.rows[i]
}

// Len returns the table's number of rows.
func (x Table) Len() int {
	return len(x.rows)
}

// Expand inserts a new row at the bottom of the table and returns it.
func (x *Table) Expand() *Row {
	x.rows = append(x.rows, Row{
		Base:  newBase(x.Call("insertRow", -1)),
		cells: make([]Cell, 0),
	})
	return &x.rows[len(x.rows)-1]
}

// Clear deletes all rows from the table.
func (x *Table) Clear() {
	for i := 0; i < len(x.rows); i++ {
		x.Base.Value.Call("deleteRow", -1)
	}
	if len(x.rows) > 0 {
		x.rows = x.rows[:0]
	}
}
