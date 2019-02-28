package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSharedMemory(t *testing.T) {
	require := require.New(t)

	id := "shm-test-" + time.Now().Format(time.RFC3339)
	size := 1234

	shm, err := NewSharedMemory(id, size)
	require.NoError(err, "can't create shm")
	defer shm.Close(true)
	require.Equal(id, shm.ID(), "bad id")
	require.Equal(size, len(shm.Data()), "bad size")
}
