package core

import (
	"os"
	"strings"
	"syscall/js"
)

type PageRoute int

const (
	ProtocolViewRoute = PageRoute(iota)
	ServicesViewRoute
	MessagesViewRoute
	ConfigViewRoute
	NotFoundRoute
)

var basePath = ""

var routes = map[string]PageRoute{
	"/":         ProtocolViewRoute,
	"/services": ServicesViewRoute,
	"/messages": MessagesViewRoute,
	"/config":   ConfigViewRoute,
}

var routesPaths map[PageRoute]string
var httpUrlBase, wsUrl string

func init() {
	if len(os.Args) >= 2 {
		basePath = strings.TrimSuffix(os.Args[1], "/")
	}

	routesPaths = map[PageRoute]string{}
	for routePath, route := range routes {
		routesPaths[route] = routePath
	}

	js.Global().Get("window").Call("addEventListener", "popstate", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		Dispatch(GetRoute())
		return nil
	}))

	// build url base
	location := js.Global().Get("window").Get("location")
	host := location.Get("host").String()
	httpUrlBase = location.Get("protocol").String()
	if httpUrlBase == "https:" {
		wsUrl = "wss://"
	} else {
		wsUrl = "ws://"
	}
	httpUrlBase += "//" + host
	wsUrl += host + "/socket"
}

func HTTPUrl(path string) string {
	return httpUrlBase + path
}

func WebSocketUrl() string {
	return wsUrl
}

func GetLocation() (path string) {
	location := js.Global().Get("window").Get("location")
	path = location.Get("pathname").String()
	return
}

func SetLocation(url string) {
	js.Global().Get("history").Call("pushState", nil, "", url)
}

func GetRoute() PageRoute {
	path := GetLocation()
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

func ToRoute(route PageRoute) bool {
	if routePath, ok := routesPaths[route]; ok {
		SetLocation(basePath + routePath)
		Dispatch(route)
		return true
	}
	return false
}
