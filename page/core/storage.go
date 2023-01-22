package core

import "syscall/js"

func LocalStorageSet(key, value string) {
	js.Global().Get("window").Get("localStorage").Call("setItem", key, value)
}

func LocalStorageGet(key string) (string, bool) {
	value := js.Global().Get("window").Get("localStorage").Call("getItem", key)
	if value.Type() != js.TypeString {
		return "", false
	}
	return value.String(), true
}

func LocalStorageRemove(key string) {
	js.Global().Get("window").Get("localStorage").Call("removeItem", key)
}

func LocalStorageClear() {
	js.Global().Get("window").Get("localStorage").Call("clear")
}

func SessionStorageSet(key, value string) {
	js.Global().Get("window").Get("sessionStorage").Call("setItem", key, value)
}

func SessionStorageGet(key string) (string, bool) {
	value := js.Global().Get("window").Get("sessionStorage").Call("getItem", key)
	if value.Type() != js.TypeString {
		return "", false
	}
	return value.String(), true
}

func SessionStorageRemove(key string) {
	js.Global().Get("window").Get("sessionStorage").Call("removeItem", key)
}

func SessionStorageClear() {
	js.Global().Get("window").Get("sessionStorage").Call("clear")
}
