// Package dom wraps JS DOM functionality.
package html

import (
	"github.com/blitz-frost/wasm"
	"github.com/blitz-frost/wasm/css"
)

type Collection struct {
	wasm.Object
}

func (x Collection) Get(idOrName string) (Element, bool) {
	v := x.CallRaw("namedItem", idOrName)
	if v.IsNull() {
		return Element{}, false
	}
	return Element{wasm.Object(v)}, true
}

func (x Collection) Index(i int) Element {
	v := x.Value().Index(i)
	return Element{wasm.Object(v)}
}

func (x Collection) Len() int {
	return x.Value().Length()
}

type ElementKind string

type Element struct {
	wasm.Object
}

// Append adds the given element as final subelement.
func (x Element) Append(elem Element) {
	x.CallRaw("appendChild", elem.Value())
}

func (x Element) Attribute(name string) (string, bool) {
	v := x.CallRaw("getAttribute", name)
	if v.IsNull() {
		return "", false
	}
	return v.String(), true
}

func (x Element) AttributeSet(name, value string) {
	x.CallRaw("setAttribute", name, value)
}

func (x Element) Blur() {
	x.CallRaw("blur")
}

func (x Element) Children() Collection {
	v := x.Get("children")
	return Collection{wasm.Object(v)}
}

func (x Element) ChildrenByKind(kind ElementKind) Collection {
	v := x.CallRaw("getElementsByTagName", string(kind))
	return Collection{wasm.Object(v)}
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
	x.CallRaw("focus")
}

// Handle subscribes the given HandlerFunction to the specified event.
func (x Element) Handle(event EventName, h HandlerFunction) {
	x.CallRaw("addEventListener", string(event), wasm.Value(h))
}

// HandleRemove unsubscribes the given Handler from the specified event.
func (x Element) HandleRemove(event EventName, h HandlerFunction) {
	x.CallRaw("removeEventListener", string(event), wasm.Value(h))
}

func (x Element) Height() uint16 {
	v := x.Get("offsetHeight")
	return uint16(v.Int())
}

func (x Element) HeightVisible() uint16 {
	v := x.Get("clientHeight")
	return uint16(v.Int())
}

func (x Element) Id() string {
	v := x.Get("id")
	return v.String()
}

func (x Element) IdSet(id string) {
	x.Set("id", id)
}

// InsertSelf inserts target element before the specified element.
func (x Element) InsertSelf(elem Element) {
	x.CallRaw("insertBefore", elem.Value())
}

func (x Element) Kind() ElementKind {
	v := x.Get("tagName")
	return ElementKind(v.String())
}

// Next returns the next element in the same node.
func (x Element) Next() (Element, bool) {
	v := x.Get("nextElementSibling")
	if v.IsNull() {
		return Element{}, false
	}
	return Element{wasm.Object(v)}, true
}

// Previous returns the previous element in the same node.
func (x Element) Previous() (Element, bool) {
	v := x.Get("previousElementSibling")
	if v.IsNull() {
		return Element{}, false
	}
	return Element{wasm.Object(v)}, true
}

// Remove removes the specified subelement.
func (x Element) Remove(elem Element) {
	x.CallRaw("removeChild", elem.Value())
}

// RemoveSelf removes the target Element from the dom
func (x Element) RemoveSelf() {
	x.CallRaw("remove")
}

func (x Element) Replace(newElem, oldElem Element) {
	x.CallRaw("replaceChild", newElem.Value(), oldElem.Value())
}

func (x Element) SpellcheckSet(val bool) {
	x.Set("spellcheck", val)
}

func (x Element) Style() css.Style {
	v := x.Get("style")
	return css.Style{wasm.Object(v)}
}

func (x Element) Parent() Element {
	v := x.Get("parentElement")
	return Element{wasm.Object(v)}
}

func (x Element) TabIndexSet(i int) {
	x.Set("tabIndex", i)
}

// Text returns the inner HTML text node value. Panics if x does not contain a text node.
func (x Element) Text() string {
	v := x.Get("innerHTML")
	return v.String()
}

// TextSet sets the inner HTML of x as a text node with the provided value.
func (x Element) TextSet(s string) {
	x.Set("innerHTML", s)
}

func (x Element) Width() uint16 {
	v := x.Get("offsetWidth")
	return uint16(v.Int())
}

func (x Element) WidthVisible() uint16 {
	v := x.Get("clientWidth")
	return uint16(v.Int())
}
