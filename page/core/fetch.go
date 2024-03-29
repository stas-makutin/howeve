package core

import (
	"syscall/js"
)

type FetchInit struct {
	Method string
	Body   string
}

type FetchResponse struct {
	OK     bool
	Status int
	Body   string
}

type FetchError struct {
	Name    string
	Message string
}

func (fe *FetchError) Error() string {
	return fe.Message
}

func Fetch(url string, init *FetchInit, then func(response *FetchResponse), catch func(err *FetchError)) {
	opts := map[string]interface{}{}
	if init != nil {
		if init.Method != "" {
			opts["method"] = init.Method
		}
		if init.Body != "" {
			opts["body"] = init.Body
		}
	}

	var thenFn, catchFn js.Func
	releaseFn := func() {
		thenFn.Release()
		catchFn.Release()
	}

	thenFn = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		releaseFn()
		r := args[0]
		response := &FetchResponse{
			OK:     r.Get("ok").Bool(),
			Status: r.Get("status").Int(),
		}
		r.Call("text").Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			response.Body = args[0].String()
			then(response)
			return nil
		}))
		return nil
	})
	catchFn = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		releaseFn()
		e := args[0]
		err := &FetchError{
			Name:    e.Get("name").String(),
			Message: e.Get("message").String(),
		}
		catch(err)
		return nil
	})

	js.Global().Call("fetch", url, opts).Call("then", thenFn).Call("catch", catchFn)
}
