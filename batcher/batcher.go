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
	input      gt.Chan[T]
	output     gt.Chan[[]T]
	stopped    notification.Notification
	bufferSize int
	timeout    time.Duration
}

func New[T any](bufferSize int, timeout time.Duration) Batcher[T] {
	batcher := &batcher[T]{
		input:      ch.New[T](),
		output:     ch.New[[]T](),
		stopped:    notification.New(),
		bufferSize: bufferSize,
		timeout:    timeout,
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
	defer ticker.Stop()

	buffer := make([]T, b.bufferSize)
	current := 0

	flush := func() {
		if current == 0 {
			return
		}
		batch := make([]T, current)
		copy(batch, buffer[0:current])
		b.output.Send(batch)
		current = 0
	}

	for {
		select {
		case msg, ok := <-b.input.Subscribe():
			if !ok {
				flush()
				b.stopped.Signal()
				return
			}
			buffer[current] = msg
			current++

			if current == len(buffer) {
				flush()
				ticker.Reset(b.timeout)
			}
		case <-ticker.C:
			flush()
		}
	}
}
