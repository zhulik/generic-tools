package batcher_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	gt "github.com/zhulik/generic-tools"
	"github.com/zhulik/generic-tools/batcher"
	"github.com/zhulik/generic-tools/notification"
)

func MustReceive[T any](r gt.Receiver[T]) T {
	v, ok := r.Receive()
	if !ok {
		panic("Not ok!")
	}
	return v
}

func TestBatcher(t *testing.T) {
	t.Parallel()
	t.Run("when sending faster than timeout", func(t *testing.T) {
		batcher := batcher.New[int](5, 2*time.Second)

		done := notification.New()

		go func() {
			assert.Equal(t, []int{0, 1, 2, 3, 4}, MustReceive[[]int](batcher))
			assert.Equal(t, []int{5, 6, 7, 8, 9}, MustReceive[[]int](batcher))
			assert.Equal(t, []int{10, 11}, MustReceive[[]int](batcher))
			done.Signal()
		}()

		for i := 0; i < 12; i++ {
			batcher.Send(i)
		}
		batcher.Close()

		done.Wait()
	})

	t.Run("when sending slower than timeout", func(t *testing.T) {
		batcher := batcher.New[int](5, 1*time.Second)

		done := notification.New()

		go func() {
			assert.Equal(t, []int{0, 1, 2}, MustReceive[[]int](batcher))
			assert.Equal(t, []int{3, 4, 5}, MustReceive[[]int](batcher))
			done.Signal()
		}()

		for i := 0; i < 6; i++ {
			batcher.Send(i)
			if i == 2 {
				time.Sleep(1100 * time.Millisecond)
			}
		}

		done.Wait()

		batcher.Close()
	})
}
