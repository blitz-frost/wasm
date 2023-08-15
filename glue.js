"use strict";

var goLoadDone;

// helper for invoking asynhronous Go functions
function goAsync(fn, ...args) {
	return new Promise((resolve, reject) => {
		fn(resolve, reject, ...args);
	});
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

// returns a promise that resolves when the Go code has finished loading
// must be called before loading the wasm module
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

// helper for invoking synchronous Go functions
function goSync(fn, ...args) {
	let o = fn(...args);
	if (o[1] === true) {
		throw o[0];
	}
	return o[0];
}
