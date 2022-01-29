package connect

import (
	"context"
	"time"
)

type deadline struct {
	timeout       context.Context
	timeoutCancel context.CancelFunc
}

func newDeadline() *deadline {
	return &deadline{timeout: context.Background()}
}

func (d *deadline) SetDeadline(t time.Time) error {
	if d.timeoutCancel != nil {
		d.timeoutCancel()
		d.timeoutCancel = nil
		d.timeout = context.Background()
	}
	if !t.IsZero() {
		d.timeout, d.timeoutCancel = context.WithTimeout(context.Background(), time.Until(t))
	}
	return nil
}
