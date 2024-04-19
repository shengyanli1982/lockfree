package ringbuffer

import (
	"sync/atomic"
	"unsafe"
)

const DefaultCircleBufferSize = 1024

type LockFreeRingBuffer struct {
	data     []interface{}
	capacity int64
	head     int64
	tail     int64
	count    int64
}

func New(capacity int) *LockFreeRingBuffer {
	if capacity <= 0 {
		capacity = DefaultCircleBufferSize
	}
	return &LockFreeRingBuffer{
		data:     make([]interface{}, capacity),
		capacity: int64(capacity),
		head:     0,
		tail:     0,
		count:    0,
	}
}

func (r *LockFreeRingBuffer) isEmpty() bool {
	return atomic.LoadInt64(&r.count) == 0
}

func (r *LockFreeRingBuffer) isFull() bool {
	return atomic.LoadInt64(&r.count) == r.capacity
}

func (r *LockFreeRingBuffer) Capacity() int64 {
	return r.capacity
}

func (r *LockFreeRingBuffer) Count() int64 {
	return r.count
}

func (r *LockFreeRingBuffer) Push(value interface{}) bool {
	for {
		if r.isFull() {
			return false
		}
		tail := atomic.LoadInt64(&r.tail)
		next := (tail + 1) % r.capacity
		if atomic.CompareAndSwapInt64(&r.tail, tail, next) {
			atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&r.data[tail])), unsafe.Pointer(&value))
			atomic.AddInt64(&r.count, 1)
			return true
		}
	}
}

func (r *LockFreeRingBuffer) Pop() (interface{}, bool) {
	for {
		if r.isEmpty() {
			return nil, false
		}

		head := atomic.LoadInt64(&r.head)
		next := (head + 1) % r.capacity
		value := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&r.data[head])))

		if atomic.CompareAndSwapInt64(&r.head, head, next) {
			atomic.AddInt64(&r.count, -1)
			return *(*interface{})(value), true
		}

	}
}

func (r *LockFreeRingBuffer) Reset() {
	atomic.StoreInt64(&r.head, 0)
	atomic.StoreInt64(&r.tail, 0)
	atomic.StoreInt64(&r.count, 0)
}
