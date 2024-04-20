package ringbuffer

// RingBuffer 是一个接口，定义了环形缓冲区的基本操作
// RingBuffer is an interface that defines basic operations of a ring buffer
type RingBuffer = interface {
	// Push 方法用于向环形缓冲区中添加一个元素，如果缓冲区已满，返回 false
	// The Push method is used to add an element to the ring buffer, returns false if the buffer is full
	Push(value interface{}) bool

	// Pop 方法用于从环形缓冲区中移除并返回一个元素，如果缓冲区为空，返回 nil 和 false
	// The Pop method is used to remove and return an element from the ring buffer, returns nil and false if the buffer is empty
	Pop() (interface{}, bool)

	// Reset 方法用于重置/清空环形缓冲区
	// The Reset method is used to reset/clear the ring buffer
	Reset()

	// Count 方法返回环形缓冲区中的元素数量
	// The Count method returns the number of elements in the ring buffer
	Count() int64
}
