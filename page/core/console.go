package core

import "syscall/js"

var Console = newJsConsole()

type jsConsole struct {
	console js.Value
}

func newJsConsole() *jsConsole {
	return &jsConsole{
		console: js.Global().Get("console"),
	}
}

func (jc *jsConsole) Write(b []byte) (int, error) {
	jc.console.Call("log", string(b))
	return len(b), nil
}

func (jc *jsConsole) Log(v interface{}) {
	jc.console.Call("log", v)
}
