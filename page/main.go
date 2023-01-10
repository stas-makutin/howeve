package main

import (
	"syscall/js"
	"time"

	"github.com/hexops/vecty"
	"github.com/stas-makutin/howeve/page/actions"
	"github.com/stas-makutin/howeve/page/core"
)

func main() {
	vecty.SetTitle("Howeve Test Page")

	setViewport()
	addStyles()
	addScript()

	// wait for MDC script
	waitTime := 6 * time.Second
	tickTime := 100 * time.Microsecond
	for js.Global().Get("mdc").IsUndefined() {
		time.Sleep(tickTime)
		if waitTime < tickTime {
			panic("failed to initialize MDC")
		}
		waitTime -= tickTime
	}

	page := newPageMain()
	if err := vecty.RenderInto("body", page); err != nil {
		panic(err)
	}

	core.Dispatch(actions.LoadEvent(0))

	select {} // run forever
}

func setViewport() {
	meta := js.Global().Get("document").Call("createElement", "meta")
	meta.Set("name", "viewport")
	meta.Set("content", "width=device-width, initial-scale=1")
	js.Global().Get("document").Get("head").Call("appendChild", meta)
}

func addStyles() {
	vecty.AddStylesheet("./material-components-web.min.css")
	vecty.AddStylesheet("./material-icons.css")

	style := js.Global().Get("document").Call("createElement", "style")
	style.Set("innerHTML", core.Stylesheet())
	js.Global().Get("document").Get("head").Call("appendChild", style)
}

func addScript() {
	script := js.Global().Get("document").Call("createElement", "script")
	script.Set("src", "./material-components-web.min.js")
	js.Global().Get("document").Get("head").Call("appendChild", script)
}
