package ringbuffer

import (
	"errors"
	"unsafe"

	shd "github.com/shengyanli1982/lockfree/internal/shared"
)

var ErrorRingBufferIsFull = errors.New("ring buffer is full")

var ErrorPushValueIsNil = errors.New("push value is nil")

type LockFreeRingBuffer struct {
	capacity uint64
	head     unsafe.Pointer
	tail     unsafe.Pointer
	wptr     unsafe.Pointer
	rptr     unsafe.Pointer
}

func New(cap int) *LockFreeRingBuffer {
	if cap <= 0 {
		cap = int(shd.DefaultRingSize)
	}

	dummy := shd.NewNode(nil)

	rb := &LockFreeRingBuffer{
		capacity: uint64(cap),
		head:     unsafe.Pointer(dummy),
		wptr:     unsafe.Pointer(dummy),
		rptr:     unsafe.Pointer(dummy),
		tail:     unsafe.Pointer(dummy),
	}

	for i := 0; i < cap; i++ {
		node := shd.NewNode(nil)
		node.Index = uint64(i)
		node.Next = unsafe.Pointer(rb.tail)
		rb.tail = unsafe.Pointer(node)
	}

	shd.LoadNode(&rb.tail).Next = unsafe.Pointer(rb.head)

	return rb
}

func (rb *LockFreeRingBuffer) Push(value interface{}) error {
	if value == nil {
		return ErrorPushValueIsNil
	}

	if rb.isFull() {
		return ErrorRingBufferIsFull
	}

	return nil
}

func (rb *LockFreeRingBuffer) Pop() interface{} {
	if rb.isEmpty() {
		return nil
	}

	return nil
}

func (rb *LockFreeRingBuffer) Length() uint64 {
	return shd.LoadNode(&rb.wptr).Index - shd.LoadNode(&rb.rptr).Index
}

func (rb *LockFreeRingBuffer) isFull() bool {
	return rb.Length() == rb.capacity
}

func (rb *LockFreeRingBuffer) isEmpty() bool {
	return rb.Length() == 0
}

func (rb *LockFreeRingBuffer) Capacity() uint64 {
	return rb.capacity
}

func (rb *LockFreeRingBuffer) Reset() {
	for i := 0; i < int(rb.capacity); i++ {
		node := shd.LoadNode(&rb.tail)
		node.ResetValue()
	}
	rb.wptr = rb.head
	rb.rptr = rb.head
}
