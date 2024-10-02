// Package wspb provides helpers for reading and writing protobuf messages.
package wspb

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/coder/websocket"
	"google.golang.org/protobuf/proto"
)

// Read reads a protobuf message from c into v.
// It will reuse buffers in between calls to avoid allocations.
func Read(ctx context.Context, c *websocket.Conn, v proto.Message) error {
	return read(ctx, c, v)
}

func read(ctx context.Context, c *websocket.Conn, v proto.Message) (err error) {
	defer errdWrap(&err, "failed to read protobuf message")

	typ, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}

	if typ != websocket.MessageBinary {
		c.Close(websocket.StatusUnsupportedData, "expected binary message")
		return fmt.Errorf("expected binary message for protobuf but got: %v", typ)
	}

	buf := bpoolGet()
	defer bpoolPut(buf)

	_, err = buf.ReadFrom(r)
	if err != nil {
		return err
	}

	err = proto.Unmarshal(buf.Bytes(), v)
	if err != nil {
		c.Close(websocket.StatusInvalidFramePayloadData, "failed to unmarshal protobuf")
		return fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return nil
}

// Write writes the protobuf message v to c.
// It will reuse buffers in between calls to avoid allocations.
func Write(ctx context.Context, c *websocket.Conn, v proto.Message) error {
	return write(ctx, c, v)
}

func write(ctx context.Context, c *websocket.Conn, v proto.Message) (err error) {
	defer errdWrap(&err, "failed to write protobuf message")

	b := bpoolGet()
	defer func() { bpoolPut(b) }()

	buf := b.Bytes()

	buf, err = proto.MarshalOptions{}.MarshalAppend(buf, v)
	if err != nil {
		return fmt.Errorf("failed to marshal protobuf: %w", err)
	}

	*b = *bytes.NewBuffer(buf) // reset buffer to new bytes

	return c.Write(ctx, websocket.MessageBinary, buf)
}

// errdWrap wraps err with fmt.Errorf if err is non nil.
// Intended for use with defer and a named error return.
// Inspired by https://github.com/golang/go/issues/32676.
func errdWrap(err *error, f string, v ...interface{}) {
	if *err != nil {
		*err = fmt.Errorf(f+": %w", append(v, *err)...)
	}
}

var bpool sync.Pool

// Get returns a buffer from the pool or creates a new one if
// the pool is empty.
func bpoolGet() *bytes.Buffer {
	b := bpool.Get()
	if b == nil {
		return bytes.NewBuffer(nil)
	}
	return b.(*bytes.Buffer)
}

// Put returns a buffer into the pool.
func bpoolPut(b *bytes.Buffer) {
	b.Reset()
	bpool.Put(b)
}
