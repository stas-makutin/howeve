package core

import "syscall/js"

type WebSocket struct {
	jsSocket js.Value
}

func NewWebSocket(url string, messageFn func(data []byte), errorFn func(), closeFn func(), openFn func()) *WebSocket {
	return &WebSocket{
		jsSocket: js.Global().Get("WebSocket").New(url),
	}
}

func (ws *WebSocket) onMessage(fn func(data string)) {
	ws.jsSocket.Call("addEventListener", "message", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			fn(args[0].String())
		}
		return nil
	}))
}

func (ws *WebSocket) onError(fn func()) {
	ws.jsSocket.Call("addEventListener", "error", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fn()
		return nil
	}))
}

func (ws *WebSocket) onOpen(fn func()) {
	ws.jsSocket.Call("addEventListener", "open", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fn()
		return nil
	}))
}

func (ws *WebSocket) onClose(fn func()) {
	ws.jsSocket.Call("addEventListener", "close", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fn()
		return nil
	}))
}

func (ws *WebSocket) send(data string) {
	ws.jsSocket.Call("send", data)
}
