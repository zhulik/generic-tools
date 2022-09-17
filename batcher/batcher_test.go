package batcher_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zhulik/generic-tools/batcher"
)

func TestBatcher(t *testing.T) {
	t.Parallel()
	t.Run("when sending faster than timeout", func(t *testing.T) {
		batcher := batcher.New[int](5, 2*time.Second)

		done := make(chan bool)

		go func() {
			assert.Equal(t, []int{0, 1, 2, 3, 4}, batcher.Receive())
			assert.Equal(t, []int{5, 6, 7, 8, 9}, batcher.Receive())
			assert.Equal(t, []int{10, 11}, batcher.Receive())
			done <- true
		}()

		for i := 0; i < 12; i++ {
			batcher.Send(i)
		}
		batcher.Close()

		<-done
	})

	t.Run("when sending slower than timeout", func(t *testing.T) {
		batcher := batcher.New[int](5, 1*time.Second)

		done := make(chan bool)

		go func() {
			assert.Equal(t, []int{0, 1, 2}, batcher.Receive())
			assert.Equal(t, []int{3, 4, 5}, batcher.Receive())
			done <- true
		}()

		for i := 0; i < 6; i++ {
			batcher.Send(i)
			if i == 2 {
				time.Sleep(1100 * time.Millisecond)
			}
		}

		<-done

		batcher.Close()
	})
}
