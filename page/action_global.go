package main

// action event when Materical Component library is initialized
type mdcInitialized int

type mdcInitializer interface {
	mdcInitialized()
}

type change struct {
	event interface{}
}

type changeNotifier interface {
	onChange(event interface{})
}

func subscribeGlobal(r interface{}) {
	dispatcherSubscribe(func(event interface{}) {
		switch e := event.(type) {
		case mdcInitialized:
			if i, ok := r.(mdcInitializer); ok {
				i.mdcInitialized()
			}
		case change:
			if i, ok := r.(changeNotifier); ok {
				i.onChange(e.event)
			}
		}
	})
}