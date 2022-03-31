package core

import (
	"syscall/js"
)

type WebSocket struct {
	jsSocket    js.Value
	onMessageFn js.Func
	onErrorFn   js.Func
	onOpenFn    js.Func
	onCloseFn   js.Func
}

func NewWebSocket(url string) *WebSocket {
	return &WebSocket{
		jsSocket: js.Global().Get("WebSocket").New(url),
	}
}

func (ws *WebSocket) OnMessage(fn func(data string)) {
	ReleaseJSFunc(&ws.onMessageFn)
	ws.onMessageFn = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			fn(args[0].Get("data").String())
		}
		return nil
	})
	ws.jsSocket.Call("addEventListener", "message", ws.onMessageFn)
}

func (ws *WebSocket) OnError(fn func()) {
	ReleaseJSFunc(&ws.onErrorFn)
	ws.onErrorFn = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fn()
		return nil
	})
	ws.jsSocket.Call("addEventListener", "error", ws.onErrorFn)
}

func (ws *WebSocket) OnOpen(fn func()) {
	ReleaseJSFunc(&ws.onOpenFn)
	ws.onOpenFn = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fn()
		return nil
	})
	ws.jsSocket.Call("addEventListener", "open", ws.onOpenFn)
}

func (ws *WebSocket) OnClose(fn func()) {
	ReleaseJSFunc(&ws.onCloseFn)
	ws.onCloseFn = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fn()
		return nil
	})
	ws.jsSocket.Call("addEventListener", "close", ws.onCloseFn)
}

func (ws *WebSocket) Send(data string) {
	ws.jsSocket.Call("send", data)
}

func (ws *WebSocket) Close() error {
	ws.jsSocket.Call("close")
	ReleaseJSFunc(&ws.onMessageFn)
	ReleaseJSFunc(&ws.onErrorFn)
	ReleaseJSFunc(&ws.onOpenFn)
	ReleaseJSFunc(&ws.onCloseFn)
	return nil
}
