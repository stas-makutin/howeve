package core

import "syscall/js"

func SafeJSValue(v *js.Value, fn func(v *js.Value) js.Value) js.Value {
	if !(v.IsUndefined() || v.IsNull()) {
		return fn(v)
	}
	return *v
}

func SafeJSOperation(v *js.Value, fn func(v *js.Value)) {
	if !(v.IsUndefined() || v.IsNull()) {
		fn(v)
	}
}

func SafeJSDestroy(v *js.Value, fn func(v *js.Value)) {
	if !(v.IsUndefined() || v.IsNull()) {
		fn(v)
		*v = js.Undefined()
	}
}

func ReleaseJSFunc(fn *js.Func) bool {
	if !(fn.Value.IsUndefined() || fn.Value.IsNull()) {
		fn.Release()
		*fn = js.Func{}
		return true
	}
	return false
}
