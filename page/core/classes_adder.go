package core

import "github.com/hexops/vecty"

type ClassAdder struct {
	Classes []string
}

func (ca ClassAdder) AddClasses(classes ...string) {
	ca.Classes = append(ca.Classes, classes...)
}

func (ca ClassAdder) ApplyClasses() vecty.Applyer {
	return vecty.MarkupIf(len(ca.Classes) > 0, vecty.Class(ca.Classes...))
}
