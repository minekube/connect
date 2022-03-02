package grpc

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type out struct {
	data []byte
	err  error
}

func newReader(ctx context.Context, o ...out) TunnelReader {
	i := 0
	return newTunnelReader(ctx, func() ([]byte, error) {
		r := o[i]
		i++
		return r.data, r.err
	})
}

func TestTunnelReader_Read(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r := newReader(ctx,
		out{data: []byte("hello")},
		out{err: io.EOF},
	)

	b := make([]byte, 2)
	n, err := r.Read(b)
	require.NoError(t, err)
	require.Equal(t, len(b), n)
	require.Equal(t, "he", string(b))

	b = make([]byte, 3+10)
	n, err = r.Read(b)
	require.NoError(t, err)
	require.Equal(t, 3, n)
	require.Equal(t, "llo", string(b[:3]))

	b = make([]byte, 1)
	n, err = r.Read(b)
	require.NotNil(t, err)
	require.ErrorIs(t, err, io.EOF)
	require.Equal(t, []byte{0}, b)
	require.Equal(t, 0, n)

	n, err = r.Read(b)
	require.NotNil(t, err)
	require.ErrorIs(t, err, io.EOF)
	require.Equal(t, 0, n)
}

func TestTunnelReader_ReadDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	r := newReader(ctx,
		out{data: []byte("hello")},
		out{err: io.EOF},
	)

	err := r.SetDeadline(time.Now().Add(time.Second / 2))
	require.NoError(t, err)

	time.Sleep(time.Second)
	b := make([]byte, 5)
	n, err := r.Read(b)
	require.Empty(t, n)
	require.NotNil(t, err)
	require.ErrorIs(t, err, os.ErrDeadlineExceeded)
	require.Equal(t, []byte{0, 0, 0, 0, 0}, b)
	require.Equal(t, 0, n)
	time.Sleep(time.Millisecond * 50)

	// reset deadline
	err = r.SetDeadline(time.Now().Add(time.Second * 1))
	require.NoError(t, err)

	b = make([]byte, 10)
	n, err = r.Read(b)
	require.NoError(t, err)
	require.Equal(t, "hello", string(b[:5]))
	require.Equal(t, 5, n)

	time.Sleep(time.Millisecond * 1300)

	b = make([]byte, 5)
	n, err = r.Read(b)
	require.Empty(t, n)
	require.NotNil(t, err)
	require.ErrorIs(t, err, os.ErrDeadlineExceeded)
	require.Equal(t, []byte{0, 0, 0, 0, 0}, b)
	require.Equal(t, 0, n)
}
