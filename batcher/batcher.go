package batcher

import (
	"time"

	"github.com/zhulik/generic-tools/ch"
	"github.com/zhulik/generic-tools/common"
)

type Batcher[T any] interface {
	common.Closer
	common.Sender[T]
	common.BatchSender[T]
	common.Receiver[[]T]
	common.Subscriber[[]T]
}

type batcher[T any] struct {
	input   ch.Chan[T]
	output  ch.Chan[[]T]
	stopped chan bool
	buffer  []T
	current int
	timeout time.Duration
}

func New[T any](sendThreshold int, timeout time.Duration) Batcher[T] {
	batcher := &batcher[T]{
		input:   ch.New[T](),
		output:  ch.New[[]T](),
		stopped: make(chan bool),
		buffer:  make([]T, sendThreshold),
		current: 0,
		timeout: timeout,
	}

	go batcher.run()

	return batcher
}

func (b *batcher[T]) Close() {
	b.input.Close()
	<-b.stopped
	b.output.Close()
}

func (b *batcher[T]) Receive() ([]T, bool) {
	return b.output.Receive()
}

func (b *batcher[T]) Subscribe() <-chan []T {
	return b.output.Subscribe()
}

func (b *batcher[T]) Send(msg T) {
	b.input.Send(msg)
}

func (b *batcher[T]) SendBatch(msgs []T) {
	for _, msg := range msgs {
		b.Send(msg)
	}
}

func (b *batcher[T]) run() {
	ticker := time.NewTicker(b.timeout)

	for {
		select {
		case msg, ok := <-b.input.Subscribe():
			if !ok {
				ticker.Stop()
				b.flush()
				b.stopped <- true
				return
			}
			b.buffer[b.current] = msg
			b.current++

			if b.current == len(b.buffer) {
				b.flush()
				ticker.Reset(b.timeout)
			}
		case <-ticker.C:
			b.flush()
		}
	}

}

func (b *batcher[T]) flush() { // TODO: make it public?
	if b.current == 0 {
		return
	}
	batch := make([]T, b.current)
	copy(batch, b.buffer[0:b.current])
	b.output.Send(batch)
	b.current = 0
}
