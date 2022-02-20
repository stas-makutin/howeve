package main

import (
	"syscall/js"
	"time"

	"github.com/hexops/vecty"
)

var console = newJsConsole()

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

func (jc *jsConsole) log(v interface{}) {
	jc.console.Call("log", v)
}

func main() {
	// dispatcherSubscribe(onAction)

	vecty.SetTitle("Howeve Test Page")

	setViewport()
	addStyles()
	addScript()

	page := newPageMain()
	if err := vecty.RenderInto("body", page); err != nil {
		panic(err)
	}

	// wait for MDC
	waitTime := 6 * time.Second
	tickTime := 100 * time.Microsecond
	for js.Global().Get("mdc").IsUndefined() {
		time.Sleep(tickTime)
		if waitTime < tickTime {
			panic("failed to initialize MDC")
		}
		waitTime -= tickTime
	}

	dispatch(mdcInitialized(0))

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
	style.Set("innerHTML", stylesheet)
	js.Global().Get("document").Get("head").Call("appendChild", style)
}

func addScript() {
	script := js.Global().Get("document").Call("createElement", "script")
	script.Set("src", "./material-components-web.min.js")
	js.Global().Get("document").Get("head").Call("appendChild", script)
}

/*
type displayDrawer bool

type pageViewStore struct {
	drawerOpen bool
}

var store = &pageViewStore{drawerOpen: false}

func onAction(event interface{}) {
	switch e := event.(type) {
	case displayDrawer:
		store.drawerOpen = bool(e)
	default:
		return
	}
	dispatch(change{store})
}

type Header struct {
	vecty.Core
}

func newHeader() (r *Header) {
	r = &Header{}
	subscribeGlobal(r)
	return
}

func (ch *Header) mdcInitialized() {
	js.Global().Get("mdc").Get("topAppBar").Get("MDCTopAppBar").Call(
		"attachTo", js.Global().Get("document").Call("querySelector", ".mdc-top-app-bar"),
	)
}

func (ch *Header) Render() vecty.ComponentOrHTML {
	return elem.Header(
		vecty.Markup(
			vecty.Class("mdc-top-app-bar", "mdc-top-app-bar--dense"),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("mdc-top-app-bar__row"),
			),
			elem.Section(
				vecty.Markup(
					vecty.Class("mdc-top-app-bar__section", "mdc-top-app-bar__section--align-start"),
				),
				elem.Button(
					vecty.Markup(
						vecty.Class("material-icons", "mdc-top-app-bar__navigation-icon", "mdc-icon-button"),
						vecty.Attribute("aria-label", "Open navigation menu"),
						event.Click(func(e *vecty.Event) {
							dispatch(displayDrawer(true))
						}),
					),
					vecty.Text("menu"),
				),
				elem.Span(
					vecty.Markup(
						vecty.Class("mdc-top-app-bar__title"),
					),
					vecty.Text("Page title"),
				),
			),
		),
	)
}

type ModalDrawer struct {
	vecty.Core
	drawerOpen bool
}

func newModalDrawer() (r *ModalDrawer) {
	r = &ModalDrawer{}
	subscribeGlobal(r)
	return
}

func (ch *ModalDrawer) mdcInitialized() {
	jsDrawer := js.Global().Get("mdc").Get("drawer").Get("MDCDrawer").Call(
		"attachTo", js.Global().Get("document").Call("querySelector", ".mdc-drawer--modal"),
	)
	jsDrawer.Call("listen", "MDCDrawer:closed", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		dispatch(displayDrawer(false))
		return nil
	}))

	js.Global().Get("mdc").Get("ripple").Get("MDCRipple").Call(
		"attachTo", js.Global().Get("document").Call("querySelector", ".mdc-list-item__ripple"),
	)
}

func (ch *ModalDrawer) onChange(event interface{}) {
	if s, ok := event.(*pageViewStore); ok {
		ch.drawerOpen = s.drawerOpen
		vecty.Rerender(ch)
	}
}

func (ch *ModalDrawer) Render() vecty.ComponentOrHTML {
	return elem.Aside(
		vecty.Markup(
			vecty.Class("mdc-drawer", "mdc-drawer--modal"),
			vecty.MarkupIf(
				ch.drawerOpen,
				vecty.Class("mdc-drawer--open"),
			),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("mdc-drawer__content"),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("mdc-list"),
				),
				elem.Anchor(
					vecty.Markup(
						vecty.Class("mdc-list-item", "mdc-list-item--activated"),
						prop.Href("#"),
						vecty.Attribute("aria-current", "page"),
						vecty.Attribute("tabindex", "0"),
					),
					elem.Span(
						vecty.Markup(
							vecty.Class("mdc-list-item__ripple"),
						),
					),
					elem.Italic(
						vecty.Markup(
							vecty.Class("material-icons", "mdc-list-item__graphic"),
							vecty.Attribute("aria-hidden", "true"),
						),
						vecty.Text("inbox"),
					),
					elem.Span(
						vecty.Markup(
							vecty.Class("mdc-list-item__text"),
						),
						vecty.Text("Nav link 1"),
					),
				),
				elem.Anchor(
					vecty.Markup(
						vecty.Class("mdc-list-item"),
						prop.Href("#"),
					),
					elem.Span(
						vecty.Markup(
							vecty.Class("mdc-list-item__ripple"),
						),
					),
					elem.Italic(
						vecty.Markup(
							vecty.Class("material-icons", "mdc-list-item__graphic"),
							vecty.Attribute("aria-hidden", "true"),
						),
						vecty.Text("send"),
					),
					elem.Span(
						vecty.Markup(
							vecty.Class("mdc-list-item__text"),
						),
						vecty.Text("Nav link 2"),
					),
				),
				elem.Anchor(
					vecty.Markup(
						vecty.Class("mdc-list-item"),
						prop.Href("#"),
					),
					elem.Span(
						vecty.Markup(
							vecty.Class("mdc-list-item__ripple"),
						),
					),
					elem.Italic(
						vecty.Markup(
							vecty.Class("material-icons", "mdc-list-item__graphic"),
							vecty.Attribute("aria-hidden", "true"),
						),
						vecty.Text("drafts"),
					),
					elem.Span(
						vecty.Markup(
							vecty.Class("mdc-list-item__text"),
						),
						vecty.Text("Nav link 3"),
					),
				),
			),
		),
	)
}

type ModalDrawerScrim struct {
	vecty.Core
}

func (p *ModalDrawerScrim) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("mdc-drawer-scrim"),
		),
	)
}

type PageView struct {
	vecty.Core
}

func (p *PageView) Render() vecty.ComponentOrHTML {
	return elem.Body(
		newHeader(),
		newModalDrawer(),
		&ModalDrawerScrim{},
		elem.Main(
			vecty.Markup(
				vecty.Class("mdc-top-app-bar--fixed-adjust"),
			),
			vecty.Text("Hello Vecty!"),
		),
	)
}
*/
