package main

import (
	"os"
	"strings"
	"syscall/js"
)

type pageRoute int

const (
	ProtocolViewRoute = pageRoute(iota)
	ServicesViewRoute
	MessagesViewRoute
	NotFoundRoute
)

var basePath = ""

var routes = map[string]pageRoute{
	"/":         ProtocolViewRoute,
	"/services": ServicesViewRoute,
	"/messages": MessagesViewRoute,
}

var routesPaths map[pageRoute]string

func init() {
	if len(os.Args) >= 2 {
		basePath = strings.TrimSuffix(os.Args[1], "/")
	}

	routesPaths = map[pageRoute]string{}
	for routePath, route := range routes {
		routesPaths[route] = routePath
	}

	js.Global().Get("window").Call("addEventListener", "popstate", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		dispatch(getRoute())
		return nil
	}))
}

func getLocation() (path string) {
	location := js.Global().Get("window").Get("location")
	path = location.Get("pathname").String()
	// hash = location.Get("hash")
	return
}

func setLocation(url string) {
	js.Global().Get("history").Call("pushState", nil, "", url)
}

func getRoute() pageRoute {
	path := getLocation()
	if strings.HasPrefix(path, basePath) {
		path = path[len(basePath):len(path)]
		for routePath, route := range routes {
			if path == routePath {
				return route
			}
		}
	}
	return NotFoundRoute
}

func toRoute(route pageRoute) bool {
	if routePath, ok := routesPaths[route]; ok {
		setLocation(basePath + routePath)
		dispatch(route)
		return true
	}
	return false
}
