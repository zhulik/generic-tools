package common

type Closer interface {
	Close()
}

type Sender[T any] interface {
	Send(T)
}

type BatchSender[T any] interface {
	SendBatch([]T)
}

type Receiver[T any] interface {
	Receive() <-chan T
}

type BatchReceiver[T any] interface {
	ReceiveBatch() <-chan []T
}
