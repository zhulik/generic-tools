package generictools

type Closer interface {
	Close()
}

type Sender[T any] interface {
	Send(T)
	SendBatch([]T)
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
	Subscriber[T]
}
