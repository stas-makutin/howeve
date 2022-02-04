package main

// action event when Materical Component library is initialized
type mdcInitialized int

type mdcInitializer interface {
	mdcInitialized()
}

func subscribeGlobal(r interface{}) {
	dispatcherSubscribe(func(event interface{}) {
		switch event.(type) {
		case mdcInitialized:
			if i, ok := r.(mdcInitializer); ok {
				i.mdcInitialized()
			}
		}
	})
}
