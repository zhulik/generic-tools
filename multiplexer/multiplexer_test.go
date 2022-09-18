package multiplexer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zhulik/generic-tools/multiplexer"
	"github.com/zhulik/generic-tools/notification"
)

func TestMultiplexer(t *testing.T) {
	t.Parallel()

	t.Run("Receive", func(t *testing.T) {
		m := multiplexer.New[int]()
		defer m.Close()

		done := notification.New()

		go func() {
			res, ok := m.Receive()
			assert.True(t, ok)

			assert.Equal(t, 1, res)
			done.Signal()
		}()

		m.Send(1)

		done.Wait()
	})

	t.Run("Subscribe", func(t *testing.T) {
		m := multiplexer.New[int]()
		defer m.Close()

		done := notification.New()

		go func() {
			for msg := range m.Subscribe() {
				assert.Equal(t, 1, msg)
				done.Signal()
			}
		}()

		m.Send(1)

		done.Wait()
	})
}
