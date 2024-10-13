/* Do no modify this file */
"use strict";

// helper for invoking asynhronous Go functions
//
// fn - Function(resolve, reject Function(Value), args ...Value)
function goAsync(fn, ...args) {
	return new Promise((resolve, reject) => {
		fn(resolve, reject, ...args);
	});
}

// can be used to create asynchronous functions
// 
// fn - Function(resolve, reject Function(Value), args ...Value)
function goAsyncClosure(fn) {
	return (...args) => {
		return goAsync(fn, ...args);
	}
}

// used by Go to catch call exceptions
function goCatchCall(obj, method, args) {
	try {
		return [false, obj[method](...args)];
	} catch(e) {
		return [true, e];
	}
}

// used by Go to catch invoke exceptions
function goCatchInvoke(fn, args) {
	try {
		return [false, fn(...args)];
	} catch(e) {
		return [true, e];
	}
}

// used by Go to catch constructor exceptions
function goCatchNew(cls, args) {
	try {
		return [false, new cls(...args)];
	} catch(e) {
		return [true, e];
	}
}

// can be used to create closure functions
// 
// fn - Function(data Value, args ...Value)
function goClosure(fn, data) {
	return (...args) => {
		return fn(data, ...args);
	}
}

// used by Go to create exported functions
function goExport(fn) {
	return (...args) => {
		let o = fn(...args);
		if (o[1] === true) {
			throw o[0];
		}
		return o[0];
	}
}

// used by Go to resolve the goLoad promise
var goLoadDone;

// Returns a promise that resolves when the Go code calls the LoadDone() function.
function goLoad() {
	return new Promise((resolve) => {
		goLoadDone = resolve;
		const go = new Go();
		fetch("static/wasm/go.wasm").then((resp) => {
			return resp.arrayBuffer();
		}).then((buffer) => {
			return WebAssembly.instantiate(buffer, go.importObject);
		}).then((res) => {
			go.run(res.instance);
		});
	});
}
