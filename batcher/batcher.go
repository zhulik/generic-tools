package batcher

import (
	"time"

	gt "github.com/zhulik/generic-tools"
	"github.com/zhulik/generic-tools/ch"
	"github.com/zhulik/generic-tools/notification"
)

type Batcher[T any] interface {
	gt.Closer
	gt.Sender[T]
	gt.BatchSender[T]
	gt.Receiver[[]T]
	gt.Subscriber[[]T]
}

type batcher[T any] struct {
	input   gt.Chan[T]
	output  gt.Chan[[]T]
	stopped notification.Notification
	buffer  []T
	current int
	timeout time.Duration
}

func New[T any](sendThreshold int, timeout time.Duration) Batcher[T] {
	batcher := &batcher[T]{
		input:   ch.New[T](),
		output:  ch.New[[]T](),
		stopped: notification.New(),
		buffer:  make([]T, sendThreshold),
		current: 0,
		timeout: timeout,
	}

	go batcher.run()

	return batcher
}

func (b *batcher[T]) Close() {
	b.input.Close()
	b.stopped.Wait()
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
				b.stopped.Signal()
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
