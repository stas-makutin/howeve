package core

import "github.com/hexops/vecty"

// If returns nil if cond is false, otherwise it returns the given children.
// fixes return type of vecty.If
func If(cond bool, children ...vecty.ComponentOrHTML) vecty.List {
	if cond {
		return vecty.List(children)
	}
	return nil
}
