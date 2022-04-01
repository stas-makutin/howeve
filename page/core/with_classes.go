package core

import (
	"github.com/hexops/vecty"
)

type Classes struct {
	Classes []string
}

func (c *Classes) WithClasses(classes ...string) {
	c.Classes = append(c.Classes, classes...)
}

func (c *Classes) ApplyClasses() vecty.Applyer {
	return vecty.MarkupIf(len(c.Classes) > 0, vecty.Class(c.Classes...))
}
