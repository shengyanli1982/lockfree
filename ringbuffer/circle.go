package ringbuffer

import (
	"sync/atomic"
	"unsafe"

	shd "github.com/shengyanli1982/lockfree/internal/shared"
)

// DefaultCircleBufferSize 是默认的环形缓冲区大小
// DefaultCircleBufferSize is the default size of the ring buffer
const DefaultCircleBufferSize = 1024

// LockFreeRingBuffer 是一个无锁环形缓冲区的结构体
// LockFreeRingBuffer is a structure of a lock-free ring buffer
type LockFreeRingBuffer struct {
	// data 是用于存储元素的切片
	// data is a slice used to store elements
	data []unsafe.Pointer

	// capacity 是环形缓冲区的容量
	// capacity is the capacity of the ring buffer
	capacity int64

	// head 是环形缓冲区的头部索引
	// head is the head index of the ring buffer
	head int64

	// tail 是环形缓冲区的尾部索引
	// tail is the tail index of the ring buffer
	tail int64
}

// New 是一个函数，用于创建一个新的 LockFreeRingBuffer 实例
// New is a function that creates a new instance of LockFreeRingBuffer
func New(capacity int) *LockFreeRingBuffer {
	// 如果传入的容量小于或等于 0，那么将容量设置为默认的环形缓冲区大小
	// If the passed in capacity is less than or equal to 0, then set the capacity to the default ring buffer size
	if capacity <= 0 {
		capacity = DefaultCircleBufferSize
	}

	// 创建一个新的 LockFreeRingBuffer 实例
	// Create a new instance of LockFreeRingBuffer
	rb := &LockFreeRingBuffer{
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
	}

	// 使用 for 循环初始化环形缓冲区的每个元素为一个新的节点，节点的值为 EmptyValue
	// Use a for loop to initialize each element of the ring buffer to a new node with a value of EmptyValue
	for i := 0; i < capacity; i++ {
		// 使用 shd.NewNode 函数创建一个新的节点，节点的值为 EmptyValue
		// Use the shd.NewNode function to create a new node with a value of EmptyValue
		node := shd.NewNode(shd.EmptyValue)

		// 将新节点的索引设置为 i
		// Set the index of the new node to i
		node.Index = int64(i)

		// 将环形缓冲区的第 i 个元素设置为新节点
		// Set the i-th element of the ring buffer to the new node
		rb.data[i] = unsafe.Pointer(node)
	}

	// 返回新创建的 LockFreeRingBuffer 实例
	// Return the newly created LockFreeRingBuffer instance
	return rb
}

// Capacity 是一个方法，返回环形缓冲区的容量
// Capacity is a method that returns the capacity of the ring buffer
func (r *LockFreeRingBuffer) Capacity() int64 {
	return r.capacity
}

// Count 是一个方法，返回环形缓冲区中的元素数量
// Count is a method that returns the number of elements in the ring buffer
func (r *LockFreeRingBuffer) Count() int64 {
	// 调用 currentState 方法获取环形缓冲区的当前状态，包括头部和尾部的位置，以及头部和尾部的指针
	// Call the currentState method to get the current state of the ring buffer, including the positions of the head and tail, and the pointers of the head and tail
	head, tail, _, tptr := r.currentState()

	// 调用 calcCount 方法计算并返回环形缓冲区中的元素数量
	// Call the calcCount method to calculate and return the number of elements in the ring buffer
	return r.calcCount(head, tail, tptr)
}

// calcCount 是一个方法，根据头部和尾部的位置，以及尾部的指针，计算环形缓冲区中的元素数量
// calcCount is a method that calculates the number of elements in the ring buffer based on the positions of the head and tail, and the pointer of the tail
func (r *LockFreeRingBuffer) calcCount(head, tail int64, tailptr unsafe.Pointer) int64 {
	// 如果头部和尾部的位置相同
	// If the positions of the head and tail are the same
	if tail == head {
		// 如果尾部的指针是 nil，或者尾部的下一个节点是 nil，那么环形缓冲区中的元素数量为 0
		// If the pointer of the tail is nil, or the next node of the tail is nil, then the number of elements in the ring buffer is 0
		if tailptr == unsafe.Pointer(nil) || shd.LoadNode(&tailptr).Next == nil {
			return 0
		}
		// 否则，环形缓冲区中的元素数量为环形缓冲区的容量
		// Otherwise, the number of elements in the ring buffer is the capacity of the ring buffer
		return r.capacity
	}

	// 如果尾部的位置大于头部的位置，那么环形缓冲区中的元素数量为尾部的位置减去头部的位置
	// If the position of the tail is greater than the position of the head, then the number of elements in the ring buffer is the position of the tail minus the position of the head
	if tail > head {
		return tail - head
	}

	// 否则，环形缓冲区中的元素数量为环形缓冲区的容量加上尾部的位置减去头部的位置
	// Otherwise, the number of elements in the ring buffer is the capacity of the ring buffer plus the position of the tail minus the position of the head
	return r.capacity + tail - head
}

// currentState 是一个方法，返回环形缓冲区的当前状态，包括头部和尾部的位置，以及头部和尾部的指针
// currentState is a method that returns the current state of the ring buffer, including the positions of the head and tail, and the pointers of the head and tail
func (r *LockFreeRingBuffer) currentState() (head, tail int64, headptr, tailptr unsafe.Pointer) {
	// 使用 atomic.LoadInt64 函数加载环形缓冲区的头部位置
	// Use the atomic.LoadInt64 function to load the position of the head of the ring buffer
	head = atomic.LoadInt64(&r.head)

	// 使用 atomic.LoadInt64 函数加载环形缓冲区的尾部位置
	// Use the atomic.LoadInt64 function to load the position of the tail of the ring buffer
	tail = atomic.LoadInt64(&r.tail)

	// 使用 atomic.LoadPointer 函数加载环形缓冲区的头部指针
	// Use the atomic.LoadPointer function to load the pointer of the head of the ring buffer
	headptr = atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&r.data[head])))

	// 使用 atomic.LoadPointer 函数加载环形缓冲区的尾部指针
	// Use the atomic.LoadPointer function to load the pointer of the tail of the ring buffer
	tailptr = atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&r.data[tail])))

	// 返回环形缓冲区的当前状态
	// Return the current state of the ring buffer
	return
}

// Reset 是一个方法，用于重置环形缓冲区
// Reset is a method that resets the ring buffer
func (r *LockFreeRingBuffer) Reset() {
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
			shd.ResetNodeWithoutIndex(shd.LoadNode(&ptr))
		}
	}

	// 使用 atomic.StoreInt64 函数将环形缓冲区的头部索引、尾部索引和元素数量都设置为 0
	// Use the atomic.StoreInt64 function to set the head index, tail index, and number of elements in the ring buffer all to 0
	atomic.StoreInt64(&r.head, 0)
	atomic.StoreInt64(&r.tail, 0)
}

// Push 方法用于向无锁环形缓冲区中推入一个元素
// The Push method is used to push an element into the lock-free ring buffer
func (r *LockFreeRingBuffer) Push(value interface{}) bool {
	// 使用无限循环，直到成功推入元素
	// Use an infinite loop until an element is successfully pushed
	for {
		// 调用 currentState 方法获取环形缓冲区的当前状态，包括头部和尾部的位置，以及头部和尾部的指针
		// Call the currentState method to get the current state of the ring buffer, including the positions of the head and tail, and the pointers of the head and tail
		head, tail, _, tailptr := r.currentState()

		// 如果环形缓冲区已满，返回 false
		// If the ring buffer is full, return false
		if r.calcCount(head, tail, tailptr) >= r.capacity {
			return false
		}

		// 计算下一个元素的位置
		// Calculate the position of the next element
		next := (tail + 1) % r.capacity

		// 使用 CAS 操作尝试修改尾部元素的位置, 并且尾部元素的指针不为空
		// Use CAS operation to try to modify the position of the tail element, and the pointer of the tail element is not null
		if atomic.CompareAndSwapInt64(&r.tail, tail, next) && tailptr != unsafe.Pointer(nil) {
			// 加载尾部节点
			// Load the tail node
			node := shd.LoadNode(&tailptr)

			// 将新元素的值设置为节点的值
			// Set the value of the new element as the value of the node
			node.Value = value

			// 将节点的下一个节点设置为一个空值
			// Set the next node of the node to an empty value
			node.Next = unsafe.Pointer(&shd.EmptyValue)

			// 返回 true，表示成功推入元素
			// Return true, indicating that the element was successfully pushed
			return true
		}
	}
}

// Pop 方法用于从无锁环形缓冲区中弹出一个元素
// The Pop method is used to pop an element from the lock-free ring buffer
func (r *LockFreeRingBuffer) Pop() (interface{}, bool) {
	// 使用无限循环，直到成功弹出元素
	// Use an infinite loop until an element is successfully popped
	for {
		// 调用 currentState 方法获取环形缓冲区的当前状态，包括头部和尾部的位置，以及头部和尾部的指针
		// Call the currentState method to get the current state of the ring buffer, including the positions of the head and tail, and the pointers of the head and tail
		head, tail, headptr, tailptr := r.currentState()

		// 如果环形缓冲区为空，返回 nil 和 false
		// If the ring buffer is empty, return nil and false
		if r.calcCount(head, tail, tailptr) <= 0 {
			return nil, false
		}

		// 计算下一个元素的位置
		// Calculate the position of the next element
		next := (head + 1) % r.capacity

		// 使用 CAS 操作尝试修改头部元素的位置，并且头部元素的指针不为空
		// Use CAS operation to try to modify the position of the head element, and the pointer of the head element is not null
		if atomic.CompareAndSwapInt64(&r.head, head, next) && headptr != unsafe.Pointer(nil) {
			// 加载头部节点
			// Load the head node
			node := shd.LoadNode(&headptr)

			// 获取节点的值
			// Get the value of the node
			value := node.Value

			// 重置节点，但不改变节点的索引
			// Reset the node, but do not change the index of the node
			shd.ResetNodeWithoutIndex(node)

			// 返回节点的值和 true，表示成功弹出元素
			// Return the value of the node and true, indicating that the element was successfully popped
			return value, true
		}
	}
}
