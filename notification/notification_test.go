package notification_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zhulik/generic-tools/notification"
)

func TestNotification(t *testing.T) {
	n := notification.New()

	go func() {
		start := time.Now()
		n.Wait()
		assert.GreaterOrEqual(t, 1*time.Second, time.Now().Sub(start))
	}()

	time.Sleep(1 * time.Second)
	n.Signal()
	n.Signal()
	n.Wait()
}
