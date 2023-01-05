package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/stas-makutin/howeve/page/core"
)

type KeyValueTableBuilder interface {
	AddDelimiterRow()
	AddKeyValueRow(key, value string)
}

type keyValueTableRows struct {
	rows vecty.List
}

func (r *keyValueTableRows) AddDelimiterRow() {
	r.rows = append(r.rows,
		elem.TableRow(
			elem.TableData(
				vecty.Markup(
					vecty.Attribute("colspan", "2"),
					vecty.Style("border-top", "1px solid gainsboro"),
				),
			),
		),
	)
}

func (r *keyValueTableRows) AddKeyValueRow(key, value string) {
	r.rows = append(r.rows,
		elem.TableRow(
			elem.TableData(
				vecty.Text(key),
			),
			elem.TableData(
				vecty.Markup(
					vecty.Style("white-space", "break-spaces"),
				),
				elem.Italic(
					vecty.Text(value),
				),
			),
		),
	)
}

type KeyValueTable struct {
	vecty.Core
	core.Classes
	core.Keyable
	buildFn func(builder KeyValueTableBuilder)
}

func NewKeyValueTable(buildFn func(builder KeyValueTableBuilder)) (r *KeyValueTable) {
	r = &KeyValueTable{buildFn: buildFn}
	return
}

func (ch *KeyValueTable) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *KeyValueTable) WithKey(key interface{}) *KeyValueTable {
	ch.Keyable.WithKey(key)
	return ch
}

func (ch *KeyValueTable) WithClasses(classes ...string) *KeyValueTable {
	ch.Classes.WithClasses(classes...)
	return ch
}

func (ch *KeyValueTable) renderKeyValueRow(name, value string) vecty.ComponentOrHTML {
	return elem.TableRow(
		elem.TableData(
			vecty.Text(name),
		),
		elem.TableData(
			vecty.Markup(
				vecty.Style("white-space", "break-spaces"),
			),
			elem.Italic(
				vecty.Text(value),
			),
		),
	)
}

func (ch *KeyValueTable) renderDelimiterRow() vecty.ComponentOrHTML {
	return elem.TableRow(
		elem.TableData(
			vecty.Markup(
				vecty.Attribute("colspan", "2"),
				vecty.Style("border-top", "1px solid gainsboro"),
			),
		),
	)
}

func (ch *KeyValueTable) Render() vecty.ComponentOrHTML {
	var rowsBuilder keyValueTableRows

	ch.buildFn(&rowsBuilder)

	return elem.Table(
		vecty.Markup(
			ch.ApplyClasses(),
		),
		elem.TableBody(
			rowsBuilder.rows,
		),
	)
}
