package batcher

import (
	"time"

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
	input   chan T
	output  chan []T
	buffer  []T
	current int
	timeout time.Duration
}

func New[T any](sendThreshold int, timeout time.Duration) Batcher[T] {
	batcher := &batcher[T]{
		input:   make(chan T),
		output:  make(chan []T),
		buffer:  make([]T, sendThreshold),
		current: 0,
		timeout: timeout,
	}

	go batcher.run()

	return batcher
}

func (b *batcher[T]) Close() {
	close(b.input)
	// TODO: wait til run() is done
	close(b.output)
}

func (b *batcher[T]) Receive() []T {
	return <-b.output
}

func (b *batcher[T]) Subscribe() <-chan []T {
	return b.output
}

func (b *batcher[T]) Send(msg T) {
	b.input <- msg
}

func (b *batcher[T]) SendBatch(msgs []T) {
	for _, msg := range msgs {
		b.input <- msg
	}
}

func (b *batcher[T]) run() {
	ticker := time.NewTicker(b.timeout)

	for {
		select {
		case msg := <-b.input:
			b.buffer[b.current] = msg
			b.current++

			if b.current == len(b.buffer)-1 {
				/// TODO: send batch to receiver
				b.current = 0
			}
			return
		case <-ticker.C:
			// TODO: send batch to receiver
		}
	}
	// TODO: once input is closed send the remaining messages to the receiver and stop the ticker
}
