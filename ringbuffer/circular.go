package ringbuffer

import (
	"unsafe"

	shd "github.com/shengyanli1982/lockfree/internal/shared"
)

type LockFreeRingBuffer struct {
	cap  uint64
	tail unsafe.Pointer
	head unsafe.Pointer
}

func New(cap uint64) *LockFreeRingBuffer {
	if cap == 0 {
		cap = shd.DefaultRingSize
	}
	return &LockFreeRingBuffer{
		cap:  cap,
		tail: unsafe.Pointer(shd.EmptyNode),
		head: unsafe.Pointer(shd.EmptyNode),
	}
}
