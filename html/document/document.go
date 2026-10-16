package document

import (
	"github.com/blitz-frost/wasm"
	"github.com/blitz-frost/wasm/html"
)

var (
	doc = wasm.Object(wasm.Global.Get("document"))
)

func Body() html.Element {
	v := doc.Get("body")
	return html.Element{wasm.Object(v)}
}

// ElementByFocus returns the currently focused element.
func ElementByFocus() (html.Element, bool) {
	v := doc.Get("activeElement")
	if v.IsNull() {
		return html.Element{}, false
	}
	return html.Element{wasm.Object(v)}, true
}

// ElementById returns the element with the given ID in the document.
// Returns an false if the ID doesn't exist.
func ElementById(id string) (html.Element, bool) {
	v := doc.CallRaw("getElementById", id)
	if v.IsNull() {
		return html.Element{}, false
	}
	return html.Element{wasm.Object(v)}, true
}

func ElementCreate(kind html.ElementKind) html.Element {
	v := doc.CallRaw("createElement", string(kind))
	return html.Element{wasm.Object(v)}
}

// NOTE it might be that this actually returns a general DOM element (either HTML or XML), but we currently don't have a way of expressing that.
func ElementCreateNS(namespace string, kind html.ElementKind) html.Element {
	v := doc.CallRaw("createElementNS", namespace, string(kind))
	return html.Element{wasm.Object(v)}
}

// ElementsByKind returns all elements of the specified kind (tag).
func ElementsByKind(kind html.ElementKind) html.Collection {
	v := doc.CallRaw("getElementsByTagName", string(kind))
	return html.Collection{wasm.Object(v)}
}

// Handle registers a document event listener.
func Handle(event html.EventName, h html.HandlerFunction) {
	doc.CallRaw("addEventListener", string(event), wasm.Value(h))
}

// HandleRemove deregisters a document event listener.
func HandleRemove(event html.EventName, h html.HandlerFunction) {
	doc.CallRaw("removeEventListener", string(event), wasm.Value(h))
}

// TextInsert inserts the given string at the current cursor position.
func TextInsert(str string) {
	doc.CallRaw("execCommand", "insertText", false, str)
}
