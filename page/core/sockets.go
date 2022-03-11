package core

import (
	"syscall/js"
)

type WebSocket struct {
	jsSocket js.Value
}

func NewWebSocket(url string) *WebSocket {
	return &WebSocket{
		jsSocket: js.Global().Get("WebSocket").New(url),
	}
}

func (ws *WebSocket) OnMessage(fn func(data string)) {
	ws.jsSocket.Call("addEventListener", "message", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			fn(args[0].Get("data").String())
		}
		return nil
	}))
}

func (ws *WebSocket) OnError(fn func()) {
	ws.jsSocket.Call("addEventListener", "error", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fn()
		return nil
	}))
}

func (ws *WebSocket) OnOpen(fn func()) {
	ws.jsSocket.Call("addEventListener", "open", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fn()
		return nil
	}))
}

func (ws *WebSocket) OnClose(fn func()) {
	ws.jsSocket.Call("addEventListener", "close", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fn()
		return nil
	}))
}

func (ws *WebSocket) Send(data string) {
	ws.jsSocket.Call("send", data)
}

func (ws *WebSocket) Close() {
	ws.jsSocket.Call("close")
}
