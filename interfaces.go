package generictools

type Closer interface {
	Close()
}

type Sender[T any] interface {
	Send(T)
	SendBatch([]T)
}

type Receiver[T any] interface {
	Receive() (T, bool)
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
}
