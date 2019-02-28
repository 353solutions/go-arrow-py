package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestShmAllocator(t *testing.T) {
	require := require.New(t)

	size := 1234
	id := "shm-allocator-" + time.Now().Format(time.RFC3339)
	a, err := NewShmAlloactor(id, size)
	require.NoError(err, "can't create allocator")
	defer a.Close(true)
	require.Equal(0, a.offset, "bad initial offset")

}
