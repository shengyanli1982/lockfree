package ringbuffer

import (
	"sync/atomic"
	"unsafe"
)

// DefaultCircleBufferSize 是默认的环形缓冲区大小
// DefaultCircleBufferSize is the default size of the ring buffer
const DefaultCircleBufferSize = 1024

// LockFreeRingBuffer 是一个无锁环形缓冲区的结构体
// LockFreeRingBuffer is a structure of a lock-free ring buffer
type LockFreeRingBuffer struct {
	// data 是用于存储元素的切片
	// data is a slice used to store elements
	data []interface{}

	// capacity 是环形缓冲区的容量
	// capacity is the capacity of the ring buffer
	capacity int64

	// head 是环形缓冲区的头部索引
	// head is the head index of the ring buffer
	head int64

	// tail 是环形缓冲区的尾部索引
	// tail is the tail index of the ring buffer
	tail int64

	// count 是环形缓冲区中的元素数量
	// count is the number of elements in the ring buffer
	count int64
}

// New 是一个函数，用于创建一个新的 LockFreeRingBuffer 实例
// New is a function that creates a new instance of LockFreeRingBuffer
func New(capacity int) *LockFreeRingBuffer {
	// 如果传入的容量小于或等于 0，那么将容量设置为默认的环形缓冲区大小
	// If the passed in capacity is less than or equal to 0, then set the capacity to the default ring buffer size
	if capacity <= 0 {
		capacity = DefaultCircleBufferSize
	}

	// 返回一个新的 LockFreeRingBuffer 实例
	// Returns a new instance of LockFreeRingBuffer
	return &LockFreeRingBuffer{
		// 使用 make 函数创建一个长度和容量都为 capacity 的切片
		// Create a slice with length and capacity both equal to capacity using the make function
		data: make([]interface{}, capacity),

		// 设置环形缓冲区的容量为 capacity
		// Set the capacity of the ring buffer to capacity
		capacity: int64(capacity),

		// 设置环形缓冲区的头部索引为 0
		// Set the head index of the ring buffer to 0
		head: 0,

		// 设置环形缓冲区的尾部索引为 0
		// Set the tail index of the ring buffer to 0
		tail: 0,

		// 设置环形缓冲区的元素数量为 0
		// Set the number of elements in the ring buffer to 0
		count: 0,
	}
}

// IsEmpty 是一个方法，用于检查环形缓冲区是否为空
// IsEmpty is a method that checks whether the ring buffer is empty
func (r *LockFreeRingBuffer) IsEmpty() bool {
	// 使用 atomic.LoadInt64 函数获取环形缓冲区中的元素数量，如果数量为 0，那么环形缓冲区为空
	// Use the atomic.LoadInt64 function to get the number of elements in the ring buffer, if the number is 0, then the ring buffer is empty
	return atomic.LoadInt64(&r.count) == 0
}

// IsFull 是一个方法，用于检查环形缓冲区是否已满
// IsFull is a method that checks whether the ring buffer is full
func (r *LockFreeRingBuffer) IsFull() bool {
	// 使用 atomic.LoadInt64 函数获取环形缓冲区中的元素数量，如果数量等于环形缓冲区的容量，那么环形缓冲区已满
	// Use the atomic.LoadInt64 function to get the number of elements in the ring buffer, if the number equals the capacity of the ring buffer, then the ring buffer is full
	return atomic.LoadInt64(&r.count) == r.capacity
}

// Capacity 是一个方法，返回环形缓冲区的容量
// Capacity is a method that returns the capacity of the ring buffer
func (r *LockFreeRingBuffer) Capacity() int64 {
	return r.capacity
}

// Count 是一个方法，返回环形缓冲区中的元素数量
// Count is a method that returns the number of elements in the ring buffer
func (r *LockFreeRingBuffer) Count() int64 {
	return r.count
}

// Reset 是一个方法，用于重置环形缓冲区
// Reset is a method that resets the ring buffer
func (r *LockFreeRingBuffer) Reset() {
	// 使用 atomic.StoreInt64 函数将环形缓冲区的头部索引、尾部索引和元素数量都设置为 0
	// Use the atomic.StoreInt64 function to set the head index, tail index, and number of elements in the ring buffer all to 0
	atomic.StoreInt64(&r.head, 0)
	atomic.StoreInt64(&r.tail, 0)
	atomic.StoreInt64(&r.count, 0)
}

func (r *LockFreeRingBuffer) Push(value interface{}) bool {
	for {
		if r.IsFull() {
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
		if r.IsEmpty() {
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
