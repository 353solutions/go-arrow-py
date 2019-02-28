package main

import (
	"fmt"
	"os"
	"syscall"
)

// TODO: Get from environment?
const (
	rootDir = "/dev/shm"
)

// SharedMemory for the poor
type SharedMemory struct {
	id   string
	data []byte
}

// NewSharedMemory creates a new shared memory
func NewSharedMemory(id string, size int) (*SharedMemory, error) {
	path := idPath(id)
	if isFile(path) {
		return nil, fmt.Errorf("shared memory %q alredy exists", id)
	}

	file, err := createSizedFile(path, size)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	data, err := syscall.Mmap(
		int(file.Fd()),
		0,
		size,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED,
	)
	if err != nil {
		return nil, err
	}

	return &SharedMemory{id, data}, nil
}

// Close closes the shared memory
func (s *SharedMemory) Close(del bool) error {
	if s.data != nil {
		syscall.Munmap(s.data)
		s.data = nil
	}

	if del {
		return os.Remove(idPath(s.id))
	}

	return nil
}

// Data return the unerlying byte array
func (s *SharedMemory) Data() []byte {
	return s.data
}

// ID returns the shared memory id
func (s *SharedMemory) ID() string {
	return s.id
}

func createSizedFile(path string, size int) (*os.File, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	ok := false
	defer func() {
		if !ok {
			file.Close()
			os.Remove(path)
		}

	}()

	// Make sure file is in given size
	if _, err = file.Seek(int64(size-1), os.SEEK_SET); err != nil {
		return nil, err
	}

	if _, err = file.Write([]byte{0}); err != nil {
		return nil, err
	}

	if _, err = file.Seek(0, os.SEEK_SET); err != nil {
		return nil, err
	}

	ok = true
	return file, nil
}

func idPath(id string) string {
	return rootDir + "/" + id
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}
