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

func (c *channel[T]) Receive() (T, bool) {
	res, ok := <-c.channel
	return res, ok
}

func (c *channel[T]) Send(msg T) {
	c.channel <- msg
}

func (c *channel[T]) Subscribe() <-chan T {
	return c.channel
}
