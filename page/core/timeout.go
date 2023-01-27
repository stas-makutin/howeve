package core

import "syscall/js"

type Timeout struct {
	ID uint
	fn js.Func
}

func (t *Timeout) Set(fn func(), delay uint) *Timeout {
	t.Clear()
	t.fn = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fn()
		return nil
	})
	t.ID = uint(js.Global().Call("setTimeout", t.fn, int(delay)).Int())
	return t
}

func (t *Timeout) Reset(delay uint) *Timeout {
	if t.ID > 0 {
		js.Global().Call("clearTimeout", int(t.ID))
		t.ID = 0
	}
	if !(t.fn.Value.IsUndefined() || t.fn.Value.IsNull()) {
		t.ID = uint(js.Global().Call("setTimeout", t.fn, int(delay)).Int())
	}
	return t
}

func (t *Timeout) Clear() {
	if t.ID > 0 {
		js.Global().Call("clearTimeout", int(t.ID))
		t.ID = 0
	}
	ReleaseJSFunc(&t.fn)
}
