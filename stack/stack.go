package stack

import (
	"unsafe"

	shd "github.com/shengyanli1982/lockfree/internal/shared"
)

// LockFreeStack 是一个无锁栈的结构体
// LockFreeStack is a structure of a lock-free stack
type LockFreeStack struct {
	// top 是栈顶元素的指针
	// top is a pointer to the top element of the stack
	top unsafe.Pointer
}

// New 函数用于创建一个新的无锁栈
// The New function is used to create a new lock-free stack
func New() *LockFreeStack {
	// 使用 shd.NewNode 函数创建一个初始化节点，将该节点作为栈的底部元素
	// Use the shd.NewNode function to create an initialization node, and use this node as the bottom element of the stack
	initNode := shd.NewNode(shd.EmptyValue)

	// 返回一个新的 LockFreeStack 实例，其中栈顶元素为刚刚创建的空节点
	// Return a new LockFreeStack instance with the top element as the newly created empty node
	return &LockFreeStack{
		top: unsafe.Pointer(initNode),
	}
}

// Push 方法用于向无锁栈中推入一个元素
// The Push method is used to push an element into the lock-free stack
func (s *LockFreeStack) Push(value interface{}) {
	// 检查值是否为空, 如果为空则直接返回
	// Check if the value is nil, if it is, return directly
	if value == nil {
		return
	}

	// 创建一个新的节点
	// Create a new node
	node := shd.NewNode(value)

	// 使用无限循环，直到成功推入元素
	// Use an infinite loop until an element is successfully pushed
	for {
		// 获取栈顶元素
		// Get the top element of the stack
		top := shd.LoadNode(&s.top)

		// 设置新节点的下一个元素为当前的栈顶元素
		// Set the next element of the new node to the current top element
		node.Next = unsafe.Pointer(top)

		// 设置新节点的索引为栈顶元素的索引加 1
		// Set the index of the new node to the index of the top element plus 1
		node.Index = top.Index + 1

		// 使用 CAS 操作尝试修改栈顶元素
		// Use CAS operation to try to modify the top element
		if shd.CompareAndSwapNode(&s.top, top, node) {
			// 结束循环
			// End the loop
			return
		}
	}
}

// Pop 方法用于从无锁栈中弹出一个元素
// The Pop method is used to pop an element from the lock-free stack
func (s *LockFreeStack) Pop() interface{} {
	// 使用无限循环，直到成功弹出元素
	// Use an infinite loop until an element is successfully popped
	for {
		// 获取栈顶元素
		// Get the top element of the stack
		top := shd.LoadNode(&s.top)

		// 获取栈顶元素的下一个元素
		// Get the next element of the top element
		next := shd.LoadNode(&top.Next)

		// 检查栈顶元素是否被其他线程修改
		// Check if the top element has been modified by other threads
		if top == shd.LoadNode(&s.top) {
			// 如果栈为空，返回 nil
			// If the stack is empty, return nil
			if next == nil {
				return nil
			}

			// 获取要返回的结果
			// Get the result to be returned
			result := top.Value

			// 使用 CAS 操作尝试修改栈顶元素
			// Use CAS operation to try to modify the top element
			if shd.CompareAndSwapNode(&s.top, top, next) {
				// 重置原栈顶元素
				// Reset the original top element
				shd.ResetNodeAll(top)

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
func (s *LockFreeStack) Length() int64 {
	// 使用 atomic.Loadint64 函数获取队列的长度
	// Use the atomic.Loadint64 function to get the length of the queue
	return shd.LoadNode(&s.top).Index
}

// Reset 方法用于重置 LockFreeQueue 队列
// The Reset method is used to reset the LockFreeQueue queue
func (s *LockFreeStack) Reset() {
	// 使用 shd.NewNode 函数创建一个初始化节点，将该节点作为栈的底部元素
	// Use the shd.NewNode function to create an initialization node, and use this node as the bottom element of the stack
	initNode := shd.NewNode(shd.EmptyValue)

	// 将队列的头节点和尾节点都设置为新创建的节点
	// Set both the head node and the tail node of the queue to the newly created node
	s.top = unsafe.Pointer(initNode)
}

// IsEmpty 是一个方法，用于检查无锁栈是否为空
// IsEmpty is a method used to check if the lock-free stack is empty
func (s *LockFreeStack) IsEmpty() bool {
	// 使用 shd.LoadNode 函数加载栈顶节点, 然后检查节点的 Index 是否为 0，如果为 0，表示栈为空
	// Use the shd.LoadNode function to load the top node, then check if the Index of the node is 0, if it is 0, the stack is empty
	return shd.LoadNode(&s.top).Index == 0
}
