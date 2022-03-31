package components

import (
	"syscall/js"

	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/hexops/vecty/prop"
	"github.com/stas-makutin/howeve/page/core"
)

// some common actions
const (
	MdcDialogActionNone    = "none"
	MdcDialogActionClose   = "close"
	MdcDialogActionOK      = "ok"
	MdcDialogActionCancel  = "cancel"
	MdcDialogActionAccept  = "accept"
	MdcDialogActionDiscard = "discard"
)

type MdcDialogButton struct {
	Label    string `vecty:"prop"`
	Action   string `vecty:"prop"`
	Disabled bool   `vecty:"prop"`
	Default  bool   `vecty:"prop"`
}

type MdcDialog struct {
	vecty.Core
	ID          string
	Title       string
	FullScreen  bool
	CloseButton bool
	Buttons     []MdcDialogButton `vecty:"prop"`
	Content     vecty.List        `vecty:"prop"`
	closeFn     func(action string)
	jsObject    js.Value
	jsClosedFn  js.Func
}

func NewMdcDialog(id, title string, fullScreen, closeButton bool, closeFn func(action string), buttons []MdcDialogButton, content ...vecty.ComponentOrHTML) (r *MdcDialog) {
	r = &MdcDialog{ID: id, Title: title, FullScreen: fullScreen, CloseButton: closeButton, Buttons: buttons, Content: content, closeFn: closeFn}
	return
}

func (ch *MdcDialog) Mount() {
	ch.Unmount()
	ch.jsObject = js.Global().Get("mdc").Get("dialog").Get("MDCDialog").Call(
		"attachTo", js.Global().Get("document").Call("getElementById", ch.ID),
	)
	core.ReleaseJSFunc(&ch.jsClosedFn)
	ch.jsClosedFn = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			ch.closeFn(args[0].Get("detail").Get("action").String())
		}
		return nil
	})
	ch.jsObject.Call("listen", "MDCDialog:closed", ch.jsClosedFn)
	ch.jsObject.Call("open")
}

func (ch *MdcDialog) Unmount() {
	core.SafeJSDestroy(&ch.jsObject, func(v *js.Value) { v.Call("destroy") })
	core.ReleaseJSFunc(&ch.jsClosedFn)
}

func (ch *MdcDialog) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *MdcDialog) Render() vecty.ComponentOrHTML {
	labelID := ch.ID + "---label"
	contentID := ch.ID + "---content"
	return elem.Div(
		vecty.Markup(
			prop.ID(ch.ID),
			vecty.Class("mdc-dialog"),
			vecty.MarkupIf(ch.FullScreen, vecty.Class("mdc-dialog--fullscreen")),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("mdc-dialog__container"),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("mdc-dialog__surface"),
					vecty.Attribute("role", "dialog"),
					vecty.Attribute("aria-modal", "true"),
					vecty.Attribute("aria-labelledby", labelID),
					vecty.Attribute("aria-describedby", contentID),
				),
				ch.RenderHeader(labelID),
				elem.Div(
					vecty.Markup(
						prop.ID(contentID),
						vecty.Class("mdc-dialog__content"),
					),
					ch.Content,
				),
				ch.RenderButtons(),
			),
		),
		elem.Div(
			vecty.Markup(
				vecty.Class("mdc-dialog__scrim"),
			),
		),
	)
}

func (ch *MdcDialog) RenderHeader(labelID string) vecty.ComponentOrHTML {
	var title *vecty.HTML
	if ch.Title != "" {
		title = elem.Heading2(
			vecty.Markup(
				prop.ID(labelID),
				vecty.Class("mdc-dialog__title"),
			),
			vecty.Text(ch.Title),
		)
	}

	if ch.CloseButton {
		return elem.Div(
			vecty.Markup(
				vecty.Class("mdc-dialog__header"),
			),
			title,
			elem.Button(
				vecty.Markup(
					vecty.Class("mdc-icon-button", "material-icons", "mdc-dialog__close"),
					vecty.Attribute("data-mdc-dialog-action", "close"),
				),
				vecty.Text("close"),
			),
		)
	}

	return title
}

func (ch *MdcDialog) RenderButtons() vecty.ComponentOrHTML {
	if len(ch.Buttons) <= 0 {
		return nil
	}

	var buttons vecty.List
	for _, btn := range ch.Buttons {
		buttons = append(buttons, elem.Button(
			vecty.Markup(
				vecty.Class("mdc-button", "mdc-dialog__button"),
				vecty.Attribute("data-mdc-dialog-action", btn.Action),
				vecty.MarkupIf(btn.Default, vecty.Attribute("data-mdc-dialog-button-default", "true")),
				prop.Disabled(btn.Disabled),
			),
			elem.Div(
				vecty.Markup(
					vecty.Class("mdc-button__ripple"),
				),
			),
			elem.Span(
				vecty.Markup(
					vecty.Class("mdc-button__label"),
				),
				vecty.Text(btn.Label),
			),
		))
	}

	return elem.Div(
		vecty.Markup(
			vecty.Class("mdc-dialog__actions"),
		),
		buttons,
	)
}
