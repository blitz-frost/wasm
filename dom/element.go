// Package dom wraps JS DOM functionality.
package dom

import (
	"syscall/js"

	"github.com/blitz-frost/wasm/css"
)

type ElementKind string

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

func (x Element) Kind() ElementKind {
	return ElementKind(x.Get("tagName").String())
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

func (x Element) SubKind(kind ElementKind) []Element {
	vals := x.Call("getElementsByTagName", string(kind))
	n := vals.Length()
	o := make([]Element, n)
	for i := 0; i < n; i++ {
		o[i] = Element{vals.Index(i)}
	}
	return o
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
