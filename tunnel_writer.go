package connect

import (
	"context"
	"errors"
	"io"
	"os"
	"sync/atomic"
	"time"
)

// TunnelWriter is the interface for writing to a tunnel connection
// while allowing to set write deadline to unblock a blocking Write call
// and return an os.ErrDeadlineExceeded.
//
// TunnelWriter allows concurrent writes but not concurrent SetDeadline calls.
// Calling SetDeadline concurrently possibly cancels other ongoing writes.
//
// SetDeadline always returns a nil-error.
type TunnelWriter interface {
	io.Writer
	SetDeadline(t time.Time) error
}

func newTunnelWriter(
	ctx context.Context,
	writeFn writeFn,
) TunnelWriter {
	w := &tunnelWriter{
		ctx:       ctx,
		deadline:  newDeadline(),
		writeNext: make(chan *write),
	}
	go func() {
		var err error
		for {
			select {
			case req := <-w.writeNext:
				err = writeFn(req.b)
				req.b = nil
				if err != nil && errors.Is(err, io.EOF) {
					w.err.Store(err)
				}
				select {
				case req.result <- err:
				case <-w.timeout.Done(): // write retryable
				case <-w.ctx.Done():
					return // stop write loop
				}
				if errors.Is(err, io.EOF) {
					return // stop write loop
				}
			case <-ctx.Done():
				return // stop write loop
			}
		}
	}()
	return w
}

type writeFn func(b []byte) (err error)

type tunnelWriter struct {
	ctx context.Context
	*deadline
	writeNext chan *write
	err       atomic.Value
}

type write struct {
	b      []byte
	result chan error
}

func (w *tunnelWriter) Write(p []byte) (n int, err error) {
	if e := w.err.Load(); e != nil {
		return 0, e.(error)
	}
	res := make(chan error)
	// send next write
	select {
	case w.writeNext <- &write{b: p, result: res}:
	case <-w.ctx.Done():
		return 0, w.ctx.Err()
	case <-w.timeout.Done():
		return 0, os.ErrDeadlineExceeded
	}
	// receive write error
	select {
	case err = <-res:
		if err == nil {
			return len(p), nil
		}
		return 0, err
	case <-w.ctx.Done():
		return 0, w.ctx.Err()
	case <-w.timeout.Done():
		return 0, os.ErrDeadlineExceeded
	}
}
