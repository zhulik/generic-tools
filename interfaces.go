package genereictools

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
	Receive() (T, bool)
}

type Subscriber[T any] interface {
	Subscribe() <-chan T
}

type Signaller interface {
	Signal()
}

type Waiter interface {
	Wait()
}

type Chan[T any] interface {
	Sender[T]
	Closer
	Receiver[T]
	Subscriber[T]
}
