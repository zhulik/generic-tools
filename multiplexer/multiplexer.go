package multiplexer

import (
	"sync"
	"sync/atomic"

	"github.com/samber/lo"
	gt "github.com/zhulik/generic-tools"
	"github.com/zhulik/generic-tools/ch"
	"github.com/zhulik/generic-tools/notification"
)

type message[T any] struct {
	payload   T
	delivered notification.Notification
}

type subscriber[T any] struct {
	ch   gt.Chan[T]
	once bool
}

type multiplexer[T any] struct {
	input       gt.Chan[message[T]]
	subscribers []subscriber[T]
	m           sync.Mutex
	stopped     notification.Notification
	closed      atomic.Bool
}

func New[T any]() gt.Chan[T] {
	m := &multiplexer[T]{
		input:       ch.New[message[T]](),
		subscribers: []subscriber[T]{},
		stopped:     notification.New(),
		closed:      atomic.Bool{},
	}

	go m.run()

	return m
}

func (m *multiplexer[T]) Close() {
	m.input.Close()
	m.closed.Store(true)
	m.stopped.Wait()
	// No need to lock the mutex here since nobody else has access to subscribers at this point
	for _, s := range m.subscribers {
		s.ch.Close()
	}
}

func (m *multiplexer[T]) Receive() (T, bool) {
	if m.closed.Load() {
		panic("Cannot receiver from a closed multiplexer")
	}
	m.m.Lock()

	c := ch.New[T]()
	s := subscriber[T]{
		ch:   c,
		once: true,
	}
	m.subscribers = append(m.subscribers, s)
	m.m.Unlock()
	res, ok := <-c.Subscribe()
	return res, ok
}

func (m *multiplexer[T]) Send(msg T) {
	mm := message[T]{
		payload:   msg,
		delivered: notification.New(),
	}
	m.input.Send(mm)
	mm.delivered.Wait()
}

func (m *multiplexer[T]) Subscribe() <-chan T {
	if m.closed.Load() {
		panic("Cannot subscribe to a closed multiplexer")
	}

	m.m.Lock()
	defer m.m.Unlock()

	c := ch.New[T]()
	s := subscriber[T]{
		ch:   c,
		once: false,
	}
	m.subscribers = append(m.subscribers, s)
	return c.Subscribe()
}

func (m *multiplexer[T]) run() {
	for msg := range m.input.Subscribe() {
		m.m.Lock()

		for _, subscriber := range m.subscribers {
			subscriber.ch.Send(msg.payload)
		}

		m.subscribers = lo.Reject(m.subscribers, func(s subscriber[T], i int) bool {
			return s.once
		})

		m.m.Unlock()
		msg.delivered.Signal()
	}
	m.stopped.Signal()
}
