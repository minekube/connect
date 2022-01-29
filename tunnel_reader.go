package connect

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"time"
)

// TunnelReader is the interface for reading from a tunnel connection
// while allowing to set a read deadline to unblock a blocking Read call
// and return an os.ErrDeadlineExceeded.
//
// TunnelReader must not be accessed concurrently.
//
// SetDeadline always returns a nil-error.
type TunnelReader interface {
	io.Reader
	SetDeadline(t time.Time) error
}

func newTunnelReader(
	ctx context.Context,
	readFn readFn,
) TunnelReader {
	r := &tunnelReader{
		ctx:      ctx,
		deadline: newDeadline(),
		readNext: make(chan chan struct{}),
	}
	go func() {
		var data []byte
		for {
			select {
			case res := <-r.readNext:
				data, r.err = readFn()
				_, _ = r.buf.Write(data)
				select {
				case res <- struct{}{}:
				case <-r.timeout.Done(): // read retryable
				case <-ctx.Done():
					return // stop read loop
				}
				if errors.Is(r.err, io.EOF) {
					return // stop read loop
				}
			case <-ctx.Done():
				return // stop read loop
			}
		}
	}()
	return r
}

type readFn func() ([]byte, error)

type tunnelReader struct {
	ctx context.Context
	*deadline
	readNext chan chan struct{}
	buf      bytes.Buffer
	err      error
}

func (t *tunnelReader) Read(p []byte) (int, error) {
	if t.err != nil {
		return 0, t.err
	}

	res := make(chan struct{})
	for t.buf.Len() < len(p) {
		// trigger next read
		select {
		case t.readNext <- res:
		case <-t.timeout.Done():
			return 0, os.ErrDeadlineExceeded
		case <-t.ctx.Done():
			return 0, t.ctx.Err()
		}
		// wait until we can read from buf
		select {
		case <-res:
			err := t.err
			if t.err != nil {
				if !errors.Is(t.err, io.EOF) {
					t.err = nil
				}
				return 0, err
			}
		case <-t.timeout.Done():
			return 0, os.ErrDeadlineExceeded
		case <-t.ctx.Done():
			return 0, t.ctx.Err()
		}
	}
	return t.buf.Read(p)
}
