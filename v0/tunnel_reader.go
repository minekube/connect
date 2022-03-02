package v0

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
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
	readChan := make(chan interface{})
	go readLoop(ctx, readFn, readChan)
	return &tunnelReader{
		ctx:      ctx,
		deadline: newDeadline(),
		readChan: readChan,
	}
}

func readLoop(ctx context.Context, rd readFn, out chan<- interface{}) {
	var v interface{}
	var err error
	for {
		v, err = rd()
		if err != nil {
			v = err
		}
		select {
		case out <- v:
		case <-ctx.Done():
			return
		}
		if errors.Is(err, io.EOF) {
			return
		}
	}
}

type readFn func() ([]byte, error)

type tunnelReader struct {
	ctx context.Context
	*deadline
	readChan <-chan interface{}
	buf      bytes.Buffer
	err      error
}

func (t *tunnelReader) Read(p []byte) (int, error) {
	if t.err != nil {
		return 0, t.err
	}
	// empty the buf first
	if t.buf.Len() != 0 {
		return t.buf.Read(p)
	}
	// Check reader is already closed or timed out
	select {
	case <-t.timeout.Done():
		return t.timedOut()
	case <-t.ctx.Done():
		return t.ctxDone()
	default:
	}
	// receive read
	select {
	case v := <-t.readChan:
		if t.err, _ = v.(error); t.err != nil {
			return 0, t.err
		}
		b := v.([]byte)
		if len(b) > len(p) {
			// buffer last bytes for next Read
			t.buf.Write(b[len(p):])
			return copy(p, b[:len(p)]), nil
		}
		return copy(p, b), nil
	case <-t.ctx.Done():
		return t.ctxDone()
	case <-t.timeout.Done():
		return t.timedOut()
	}
}

func (t *tunnelReader) ctxDone() (int, error) {
	if errors.Is(t.ctx.Err(), context.DeadlineExceeded) {
		return 0, os.ErrDeadlineExceeded
	}
	return 0, net.ErrClosed
}

func (t *tunnelReader) timedOut() (int, error) {
	return 0, os.ErrDeadlineExceeded
}
