package connect

import (
	"context"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestTunnelWriter_Write(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	b := []byte("hello")
	var c int
	w := newTunnelWriter(ctx, func(b []byte) (err error) {
		c++
		require.Equal(t, "hello", string(b))
		return nil
	})
	n, err := w.Write(b)
	require.NoError(t, err)
	require.Equal(t, len(b), n)
	require.Equal(t, "hello", string(b))

	require.Equal(t, 1, c)

	n, err = w.Write(b)
	require.NoError(t, err)
	require.Equal(t, len(b), n)
	require.Equal(t, "hello", string(b))

	require.Equal(t, 2, c)
}

func TestTunnelWriter_WriteDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	b := []byte("hello")
	var c int
	w := newTunnelWriter(ctx, func(b []byte) (err error) {
		c++
		require.Equal(t, "hello", string(b))
		return nil
	})
	err := w.SetDeadline(time.Now().Add(time.Second / 2))
	require.NoError(t, err)

	time.Sleep(time.Second)
	n, err := w.Write(b)
	require.ErrorIs(t, err, os.ErrDeadlineExceeded)
	require.Equal(t, 0, n)
	time.Sleep(time.Millisecond * 50)

	// reset deadline
	err = w.SetDeadline(time.Now().Add(time.Second * 1))
	require.NoError(t, err)

	n, err = w.Write(b)
	require.NoError(t, err)
	require.Equal(t, 5, n)

	time.Sleep(time.Millisecond * 1300)

	n, err = w.Write(b)
	require.ErrorIs(t, err, os.ErrDeadlineExceeded)
	require.Equal(t, 0, n)
}
