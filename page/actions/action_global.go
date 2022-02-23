package actions

import "github.com/stas-makutin/howeve/page/core"

// load page event
type LoadEvent int

type LoadNotifier interface {
	OnLoad()
}

// change event
type ChangeEvent struct {
	Event interface{}
}

type ChangeNotifier interface {
	OnChange(event interface{})
}

type RouteNotifier interface {
	OnRouteChange(route core.PageRoute)
}

func SubscribeGlobal(r interface{}) {
	core.DispatcherSubscribe(func(event interface{}) {
		switch e := event.(type) {
		case LoadEvent:
			if i, ok := r.(LoadNotifier); ok {
				i.OnLoad()
			}
		case ChangeEvent:
			if i, ok := r.(ChangeNotifier); ok {
				i.OnChange(e.Event)
			}
		case core.PageRoute:
			if i, ok := r.(RouteNotifier); ok {
				i.OnRouteChange(e)
			}
		}
	})
}
