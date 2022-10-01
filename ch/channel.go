package ch

import (
	gt "github.com/zhulik/generic-tools"
)

type channel[T any] struct {
	channel chan T
}

func New[T any]() gt.Chan[T] {
	return &channel[T]{
		channel: make(chan T),
	}
}

func NewBuf[T any](size int) gt.Chan[T] {
	return &channel[T]{
		channel: make(chan T, size),
	}
}

func From[T any](c chan T) gt.Chan[T] {
	return &channel[T]{
		channel: c,
	}
}

func (c *channel[T]) Close() {
	close(c.channel)
}

func (c *channel[T]) Send(msg T) {
	c.channel <- msg
}

func (c *channel[T]) SendBatch(msgs []T) {
	for _, msg := range msgs {
		c.Send(msg)
	}
}

func (c *channel[T]) Subscribe() <-chan T {
	return c.channel
}
