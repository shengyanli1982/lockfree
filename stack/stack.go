package stack

import (
	"sync/atomic"
	"unsafe"

	shd "github.com/shengyanli1982/lockfree/internal/shared"
)

// LockFreeStack 是一个无锁栈的结构体
// LockFreeStack is a structure of a lock-free stack
type LockFreeStack struct {
	// length 是栈的长度
	// length is the length of the stack
	length uint64

	// top 是栈顶元素的指针
	// top is a pointer to the top element of the stack
	top unsafe.Pointer
}

// New 函数用于创建一个新的无锁栈
// The New function is used to create a new lock-free stack
func New() *LockFreeStack {
	// 返回一个新的 LockFreeStack，栈顶元素为空节点
	// Return a new LockFreeStack with the top element as an empty node
	return &LockFreeStack{
		top: unsafe.Pointer(shd.NewNode(shd.EmptyValue)),
	}
}

// Push 方法用于向无锁栈中推入一个元素
// The Push method is used to push an element into the lock-free stack
func (q *LockFreeStack) Push(value interface{}) {
	// 创建一个新的节点
	// Create a new node
	node := shd.NewNode(value)

	// 使用无限循环，直到成功推入元素
	// Use an infinite loop until an element is successfully pushed
	for {
		// 获取栈顶元素
		// Get the top element of the stack
		top := shd.LoadNode(&q.top)

		// 设置新节点的下一个元素为当前的栈顶元素
		// Set the next element of the new node to the current top element
		node.Next = unsafe.Pointer(top)

		// 使用 CAS 操作尝试修改栈顶元素
		// Use CAS operation to try to modify the top element
		if shd.CompareAndSwapNode(&q.top, top, node) {
			// 如果成功修改，栈的长度加 1
			// If the modification is successful, the length of the stack is increased by 1
			atomic.AddUint64(&q.length, 1)

			// 结束循环
			// End the loop
			return
		}
	}
}

// Pop 方法用于从无锁栈中弹出一个元素
// The Pop method is used to pop an element from the lock-free stack
func (q *LockFreeStack) Pop() interface{} {
	// 使用无限循环，直到成功弹出元素
	// Use an infinite loop until an element is successfully popped
	for {
		// 获取栈顶元素
		// Get the top element of the stack
		top := shd.LoadNode(&q.top)

		// 获取栈顶元素的下一个元素
		// Get the next element of the top element
		next := shd.LoadNode(&top.Next)

		// 检查栈顶元素是否被其他线程修改
		// Check if the top element has been modified by other threads
		if top == shd.LoadNode(&q.top) {

			// 如果栈为空，返回 nil
			// If the stack is empty, return nil
			if next == nil {
				return nil
			}

			// 获取要返回的结果
			// Get the result to be returned
			result := next.Value

			// 使用 CAS 操作尝试修改栈顶元素
			// Use CAS operation to try to modify the top element
			if shd.CompareAndSwapNode(&q.top, top, next) {

				// 如果成功修改，栈的长度减 1
				// If the modification is successful, the length of the stack is reduced by 1
				atomic.AddUint64(&q.length, ^uint64(0))

				// 重置原栈顶元素
				// Reset the original top element
				top.ResetAll()

				// 检查结果是否为空值
				// Check if the result is an empty value
				if result == shd.EmptyValue {
					// 如果结果是空值，返回 nil
					// If the result is an empty value, return nil
					return nil
				} else {
					// 如果结果不是空值，返回结果
					// If the result is not an empty value, return the result
					return result
				}
			}
		}
	}
}

// Length 方法用于获取 LockFreeQueue 队列的长度
// The Length method is used to get the length of the LockFreeQueue queue
func (q *LockFreeStack) Length() uint64 {
	// 使用 atomic.LoadUint64 函数获取队列的长度
	// Use the atomic.LoadUint64 function to get the length of the queue
	return atomic.LoadUint64(&q.length)
}

// Reset 方法用于重置 LockFreeQueue 队列
// The Reset method is used to reset the LockFreeQueue queue
func (q *LockFreeStack) Reset() {
	// 将队列的头节点和尾节点都设置为新创建的节点
	// Set both the head node and the tail node of the queue to the newly created node
	q.top = unsafe.Pointer(shd.NewNode(shd.EmptyValue))

	// 使用 atomic.StoreUint64 函数将队列的长度设置为 0
	// Use the atomic.StoreUint64 function to set the length of the queue to 0
	atomic.StoreUint64(&q.length, 0)
}
