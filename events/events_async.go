package events

import "sync"

type packet struct {
	event     interface{}
	receivers []SubscriberID
}

// AsyncDispatcher extends event Dispatcher to support asyncronous Send (in separate gorutine)
type AsyncDispatcher struct {
	Dispatcher
	senderCh chan struct{}
	packetCh chan *packet
	stopCh   chan struct{}
	stopWg   sync.WaitGroup
}

// NewAsyncDispatcher creates and initializes new asyncronous event Dispatcher
// The maxSenders parameter defines maximal number of goroutines used to send events (if <= 0 then 1 goroutine will be used)
func NewAsyncDispatcher(maxSenders int) *AsyncDispatcher {
	if maxSenders < 1 {
		maxSenders = 1
	}

	d := &AsyncDispatcher{
		senderCh: make(chan struct{}, maxSenders),
		packetCh: make(chan *packet, maxSenders*4),
		stopCh:   make(chan struct{}),
	}

	for i := 0; i < maxSenders; i++ {
		d.senderCh <- struct{}{}
	}

	d.stopWg.Add(1)
	go d.asyncDispatcher()

	return d
}

// Close releases resources used by asyncronous event Dispatcher
func (d *AsyncDispatcher) Close() {
	close(d.stopCh)
	d.stopWg.Wait()
	close(d.packetCh)
	close(d.senderCh)
}

// SendAsync sends event asyncronously, in a separate goroutine.
func (d *AsyncDispatcher) SendAsync(event interface{}, receivers ...SubscriberID) {
	d.packetCh <- &packet{event: event, receivers: receivers}
}

func (d *AsyncDispatcher) asyncDispatcher() {
	defer d.stopWg.Done()
	for {
		var packet *packet

		// recieve event
		select {
		case <-d.stopCh:
			return
		case packet = <-d.packetCh:
		}

		// send event
		select {
		case <-d.stopCh:
			return
		case <-d.senderCh:
		}

		go d.send(packet)
	}
}

func (d *AsyncDispatcher) send(packet *packet) {
	defer d.stopWg.Done()
	defer func() { d.senderCh <- struct{}{} }()
	d.Send(packet.event, packet.receivers...)
}
