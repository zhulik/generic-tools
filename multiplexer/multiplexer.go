package multiplexer

import (
	"sync"
	"sync/atomic"

	gt "github.com/zhulik/generic-tools"
	"github.com/zhulik/generic-tools/ch"
	"github.com/zhulik/generic-tools/notification"
)

type multiplexer[T any] struct {
	input       gt.Chan[T]
	subscribers []gt.Chan[T]
	m           sync.RWMutex
	subscribed  notification.Notification
	stopped     notification.Notification
	closed      atomic.Bool
}

func New[T any]() gt.Chan[T] {
	m := &multiplexer[T]{
		input:       ch.New[T](),
		subscribers: []gt.Chan[T]{},
		subscribed:  notification.New(),
		stopped:     notification.New(),
		closed:      atomic.Bool{},
	}

	go m.run()

	return m
}

func (m *multiplexer[T]) Close() {
	if m.closed.Load() {
		return
	}
	m.closed.Store(true)
	m.input.Close()
	m.stopped.Wait()
	// No need to lock the mutex here since nobody else has access to subscribers at this point
	for _, s := range m.subscribers {
		s.Close()
	}
}

func (m *multiplexer[T]) Send(msg T) {
	m.subscribed.Wait()
	m.input.Send(msg)
}

func (m *multiplexer[T]) SendBatch(msgs []T) {
	for _, msg := range msgs {
		m.Send(msg)
	}
}

func (m *multiplexer[T]) Subscribe() <-chan T {
	if m.closed.Load() {
		panic("Cannot subscribe to a closed multiplexer")
	}

	m.m.Lock()
	defer m.m.Unlock()

	c := ch.New[T]()
	m.subscribers = append(m.subscribers, c)
	m.subscribed.Signal()
	return c.Subscribe()
}

func (m *multiplexer[T]) run() {
	for msg := range m.input.Subscribe() {
		m.m.RLock()

		subs := make([]gt.Chan[T], len(m.subscribers))
		copy(subs, m.subscribers)

		m.m.RUnlock()

		for _, subscriber := range subs {
			subscriber.Send(msg)
		}
	}
	m.stopped.Signal()
}
