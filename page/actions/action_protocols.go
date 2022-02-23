package actions

import "github.com/stas-makutin/howeve/page/core"

// actions

type ProtocolsLoad int

const (
	ProtocolsLoadFetch = ProtocolsLoad(iota)
	ProtocolsLoadSocket
)

func protocolsLoadByFetch() {
	core.Console.Log("protocols: fetch")
}

func protocolsLoadBySocket() {
	core.Console.Log("protocols: websockets")
}

// store

type ProtocolViewStore struct {
	Loading bool
}

var pvStore = &ProtocolViewStore{
	Loading: true,
}

// reducer

func ProtocolViewRegister() {
	core.DispatcherSubscribe(pvAction)
}

func pvAction(event interface{}) {
	switch e := event.(type) {
	case ProtocolsLoad:
		pvStore.Loading = true
		if e == ProtocolsLoadFetch {
			protocolsLoadByFetch()
		} else {
			protocolsLoadBySocket()
		}
	default:
		return
	}
	core.Dispatch(ChangeEvent{pvStore})
}
