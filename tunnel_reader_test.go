package connect

import (
	"context"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestTunnelReader_Read(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var c int
	r := newTunnelReader(ctx, func() ([]byte, error) {
		if c == 1 {
			time.Sleep(time.Hour) // block
			return nil, nil
		}
		c++
		return []byte("hello"), nil
	})
	b := make([]byte, 2)
	n, err := r.Read(b)
	require.NoError(t, err)
	require.Equal(t, len(b), n)
	require.Equal(t, "he", string(b))

	b = make([]byte, 3)
	n, err = r.Read(b)
	require.NoError(t, err)
	require.Equal(t, len(b), n)
	require.Equal(t, "llo", string(b))

	b = make([]byte, 1)
	n, err = r.Read(b)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Equal(t, []byte{0}, b)
	require.Equal(t, 0, n)
}

func TestTunnelReader_ReadDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	r := newTunnelReader(ctx, func() ([]byte, error) {
		return []byte("hello"), nil
	})
	err := r.SetDeadline(time.Now().Add(time.Second / 2))
	require.NoError(t, err)

	time.Sleep(time.Second)
	b := make([]byte, 5)
	n, err := r.Read(b)
	require.ErrorIs(t, err, os.ErrDeadlineExceeded)
	require.Equal(t, []byte{0, 0, 0, 0, 0}, b)
	require.Equal(t, 0, n)
	time.Sleep(time.Millisecond * 50)

	// reset deadline
	err = r.SetDeadline(time.Now().Add(time.Second * 1))
	require.NoError(t, err)

	b = make([]byte, 5)
	n, err = r.Read(b)
	require.NoError(t, err)
	require.Equal(t, "hello", string(b))
	require.Equal(t, 5, n)

	time.Sleep(time.Millisecond * 1300)

	b = make([]byte, 5)
	n, err = r.Read(b)
	require.ErrorIs(t, err, os.ErrDeadlineExceeded)
	require.Equal(t, []byte{0, 0, 0, 0, 0}, b)
	require.Equal(t, 0, n)
}
