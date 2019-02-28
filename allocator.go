package main

import (
	"github.com/apache/arrow/go/arrow/memory"
)

const (
	memAlign = 64
)

var (
	// Make sure ShmAllocator implements memory.Allocator
	_ memory.Allocator = &ShmAllocator{}
)

// ShmAllocator is a shared memory allocator
type ShmAllocator struct {
	shm    *SharedMemory
	offset int
}

// NewShmAlloactor returns a new shared memory allocator
func NewShmAlloactor(id string, maxSize int) (*ShmAllocator, error) {
	size := align(maxSize, memAlign)
	shm, err := NewSharedMemory(id, size)
	if err != nil {
		return nil, err
	}

	return &ShmAllocator{shm: shm}, nil
}

// Allocate memory
func (a *ShmAllocator) Allocate(size int) []byte {
	size = align(size, memAlign)
	data := a.shm.Data()
	if a.offset+size > cap(data) {
		panic("out of memory")
	}
	a.offset += size
	return data[a.offset:size]
}

// Reallocate reallocates memory
func (a *ShmAllocator) Reallocate(size int, b []byte) []byte {
	if size == len(b) {
		return b
	}

	data := a.Allocate(size)
	copy(data, b)
	return data
}

// Free frees the memory
func (a *ShmAllocator) Free(b []byte) {
	// TODO: Keep free list?
}

// Close closes the shared memory
func (a *ShmAllocator) Close(del bool) error {
	if a.shm == nil {
		return nil
	}

	err := a.shm.Close(del)
	a.shm = nil
	return err
}

func align(num, size int) int {
	n := (num + size - 1) / size
	return n * size
}
