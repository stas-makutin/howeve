package components

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
	"github.com/stas-makutin/howeve/page/core"
)

type MdcGridCell struct {
	vecty.Core
	core.ClassAdder
	Content vecty.List `vecty:"prop"`
}

func NewMdcGridCell(content ...vecty.ComponentOrHTML) *MdcGridCell {
	return &MdcGridCell{Content: content}
}

func (ch *MdcGridCell) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *MdcGridCell) AddClasses(classes ...string) vecty.Component {
	ch.ClassAdder.AddClasses(classes...)
	return ch
}

func (ch *MdcGridCell) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("mdc-layout-grid__cell"),
			ch.ApplyClasses(),
		),
		ch.Content,
	)
}

type MdcGridRow struct {
	vecty.Core
	core.ClassAdder
	Cells vecty.List `vecty:"prop"`
}

func NewMdcGridRow(cells ...vecty.ComponentOrHTML) *MdcGridRow {
	return &MdcGridRow{Cells: cells}
}

func NewMdcGridSingleCellRow(context ...vecty.ComponentOrHTML) *MdcGridRow {
	cell := NewMdcGridCell(context...).AddClasses("mdc-layout-grid__cell--span-12")
	return &MdcGridRow{Cells: []vecty.ComponentOrHTML{cell}}
}

func (ch *MdcGridRow) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *MdcGridRow) AddClasses(classes ...string) vecty.Component {
	ch.ClassAdder.AddClasses(classes...)
	return ch
}

func (ch *MdcGridRow) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("mdc-layout-grid__inner"),
			ch.ApplyClasses(),
		),
		ch.Cells,
	)
}

type MdcGrid struct {
	vecty.Core
	core.ClassAdder
	Rows vecty.List `vecty:"prop"`
}

func NewMdcGrid(rows ...vecty.ComponentOrHTML) *MdcGrid {
	return &MdcGrid{Rows: rows}
}

func (ch *MdcGrid) Copy() vecty.Component {
	cpy := *ch
	return &cpy
}

func (ch *MdcGrid) AddClasses(classes ...string) vecty.Component {
	ch.ClassAdder.AddClasses(classes...)
	return ch
}

func (ch *MdcGrid) Render() vecty.ComponentOrHTML {
	return elem.Div(
		vecty.Markup(
			vecty.Class("mdc-layout-grid"),
			ch.ApplyClasses(),
		),
		ch.Rows,
	)
}
