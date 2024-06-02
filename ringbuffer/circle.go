package ringbuffer

import (
	"runtime"
	"sync/atomic"
	"unsafe"

	shd "github.com/shengyanli1982/lockfree/internal/shared"
)

// DefaultCircleBufferSize 是默认的环形缓冲区大小
// DefaultCircleBufferSize is the default size of the ring buffer
const DefaultCircleBufferSize = 1024

// rbufferImpl 是一个无锁环形缓冲区的结构体
// rbufferImpl is a structure of a lock-free ring buffer
type rbufferImpl struct {
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

	// data 是用于存储元素的切片
	// data is a slice used to store elements
	data []unsafe.Pointer
}

// New 是一个函数，用于创建一个新的 LockFreeRingBuffer 实例
// New is a function that creates a new instance of LockFreeRingBuffer
func New(capacity int) RingBuffer {
	// 如果传入的容量小于或等于 0，那么将容量设置为默认的环形缓冲区大小
	// If the passed in capacity is less than or equal to 0, then set the capacity to the default ring buffer size
	if capacity <= 0 {
		capacity = DefaultCircleBufferSize
	}

	// 创建一个新的 LockFreeRingBuffer 实例
	// Create a new instance of LockFreeRingBuffer
	rb := &rbufferImpl{
		// 使用 make 函数创建一个长度和容量都为 capacity 的切片
		// Create a slice with length and capacity both equal to capacity using the make function
		data: make([]unsafe.Pointer, capacity),

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

	// 使用 for 循环初始化环形缓冲区的每个元素为一个新的节点，节点的值为 EmptyValue
	// Use a for loop to initialize each element of the ring buffer to a new node with a value of EmptyValue
	for i := 0; i < capacity; i++ {
		rb.data[i] = unsafe.Pointer(shd.NewNode(nil))
	}

	// 返回新创建的 LockFreeRingBuffer 实例
	// Return the newly created LockFreeRingBuffer instance
	return rb
}

// IsEmpty 是一个方法，用于检查环形缓冲区是否为空
// IsEmpty is a method that checks whether the ring buffer is empty
func (r *rbufferImpl) IsEmpty() bool {
	return r.Count() == 0
}

// IsFull 是一个方法，用于检查环形缓冲区是否已满
// IsFull is a method that checks whether the ring buffer is full
func (r *rbufferImpl) IsFull() bool {
	return r.Count() == r.capacity
}

// Capacity 是一个方法，返回环形缓冲区的容量
// Capacity is a method that returns the capacity of the ring buffer
func (r *rbufferImpl) Capacity() int64 {
	return r.capacity
}

// Count 是一个方法，返回环形缓冲区中的元素数量
// Count is a method that returns the number of elements in the ring buffer
func (r *rbufferImpl) Count() int64 {
	return r.count
}

// Reset 是一个方法，用于重置环形缓冲区
// Reset is a method that resets the ring buffer
func (r *rbufferImpl) Reset() {
	// 使用 for 循环遍历环形缓冲区的每个元素
	// Use a for loop to traverse each element of the ring buffer
	for i := int64(0); i < r.capacity; i++ {
		// 使用 atomic.LoadPointer 函数获取当前元素的指针
		// Use the atomic.LoadPointer function to get the pointer of the current element
		ptr := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&r.data[i])))

		// 如果当前元素的指针不为空
		// If the pointer of the current element is not null
		if ptr != unsafe.Pointer(nil) {
			// 使用 LoadNode 方法获取节点，并调用 ResetAll 方法重置节点
			// Use the LoadNode method to get the node and call the ResetAll method to reset the node
			shd.LoadNode(&ptr).ResetAll()
		}
	}

	// 使用 atomic.StoreInt64 函数将环形缓冲区的头部索引、尾部索引和元素数量都设置为 0
	// Use the atomic.StoreInt64 function to set the head index, tail index, and number of elements in the ring buffer all to 0
	atomic.StoreInt64(&r.head, 0)
	atomic.StoreInt64(&r.tail, 0)
	atomic.StoreInt64(&r.count, 0)
}

// Push 方法用于向无锁环形缓冲区中推入一个元素
// The Push method is used to push an element into the lock-free ring buffer
func (r *rbufferImpl) Push(value interface{}) bool {
	// 使用无限循环，直到成功推入元素
	// Use an infinite loop until an element is successfully pushed
	for {
		// 如果缓冲区已满，返回 false
		// If the buffer is full, return false
		if r.IsFull() {
			return false
		}

		// 获取尾部元素的位置
		// Get the position of the tail element
		tail := atomic.LoadInt64(&r.tail)

		// 计算下一个元素的位置
		// Calculate the position of the next element
		next := (tail + 1) % r.capacity

		// 获取尾部元素的指针
		// Get the pointer of the tail element
		ptr := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&r.data[tail])))

		// 使用 CAS 操作尝试修改尾部元素的位置, 并且尾部元素的指针不为空
		// Use CAS operation to try to modify the position of the tail element, and the pointer of the tail element is not null
		if shd.CompareAndSwapInt64(&r.tail, tail, next) && ptr != unsafe.Pointer(nil) {
			// 缓冲区的元素数量加 1
			// The number of elements in the buffer is increased by 1
			atomic.AddInt64(&r.count, 1)

			// 修改尾部元素的值
			// Modify the value of the tail element
			shd.LoadNode(&ptr).Value = value

			// 返回 true，表示成功推入元素
			// Return true, indicating that the element was successfully pushed
			return true
		} else {
			// 如果 CAS 操作失败，调用 runtime.Gosched 函数让出当前线程的执行权限
			// If the CAS operation fails, call the runtime.Gosched function to yield the execution permission of the current thread
			runtime.Gosched()
		}
	}
}

// Pop 方法用于从无锁环形缓冲区中弹出一个元素
// The Pop method is used to pop an element from the lock-free ring buffer
func (r *rbufferImpl) Pop() (interface{}, bool) {
	// 使用无限循环，直到成功弹出元素
	// Use an infinite loop until an element is successfully popped
	for {
		// 如果缓冲区为空，返回 nil 和 false
		// If the buffer is empty, return nil and false
		if r.IsEmpty() {
			return nil, false
		}

		// 获取头部元素的位置
		// Get the position of the head element
		head := atomic.LoadInt64(&r.head)

		// 计算下一个元素的位置
		// Calculate the position of the next element
		next := (head + 1) % r.capacity

		// 获取头部元素的指针
		// Get the pointer of the head element
		ptr := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&r.data[head])))

		// 使用 CAS 操作尝试修改头部元素的位置，并且头部元素的指针不为空
		// Use CAS operation to try to modify the position of the head element, and the pointer of the head element is not null
		if shd.CompareAndSwapInt64(&r.head, head, next) && ptr != unsafe.Pointer(nil) {
			// 如果成功修改，缓冲区的元素数量减 1
			// If the modification is successful, the number of elements in the buffer is reduced by 1
			atomic.AddInt64(&r.count, -1)

			// 获取节点
			// Get the node
			node := shd.LoadNode(&ptr)

			// 获取节点的值
			// Get the value of the node
			value := node.Value

			// 重置节点
			// Reset the node
			node.ResetAll()

			// 返回节点的值和 true
			// Return the value of the node and true
			return value, true
		} else {
			// 如果 CAS 操作失败，调用 runtime.Gosched 函数让出当前线程的执行权限
			// If the CAS operation fails, call the runtime.Gosched function to yield the execution permission of the current thread
			runtime.Gosched()
		}
	}
}
