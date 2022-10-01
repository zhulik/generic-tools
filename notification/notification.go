package notification

import (
	"sync"
	"sync/atomic"

	gt "github.com/zhulik/generic-tools"
)

type Notification interface {
	gt.Signaller
	gt.Waiter
	Done() bool
}

type notification struct {
	wg   *sync.WaitGroup
	done atomic.Bool
}

func New() Notification {
	wg := sync.WaitGroup{}
	wg.Add(1)
	return &notification{
		wg:   &wg,
		done: atomic.Bool{},
	}
}

func (n *notification) Signal() {
	if !n.done.Swap(true) {
		n.wg.Done()
	}
}

func (n *notification) Wait() {
	if n.done.Load() {
		return
	}
	n.wg.Wait()
}

func (n *notification) Done() bool {
	return n.done.Load()
}
