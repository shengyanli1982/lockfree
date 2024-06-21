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
	length int64

	// top 是栈顶元素的指针
	// top is a pointer to the top element of the stack
	top unsafe.Pointer

	// pool 是一个节点池，用于存储和获取节点
	// pool is a node pool used to store and retrieve nodes
	pool *shd.NodePool
}

// New 函数用于创建一个新的无锁栈
// The New function is used to create a new lock-free stack
func New() *LockFreeStack {
	// 调用 newLFS 函数创建一个新的 LockFreeStack 栈，参数为 nil
	// Call the newLFS function to create a new LockFreeStack stack, the parameter is nil
	return newLFS(nil)
}

// NewWithPool 函数用于创建一个新的 LockFreeStack 栈，该栈使用一个节点池
// The NewWithPool function is used to create a new LockFreeStack stack, this stack uses a node pool
func NewWithPool() *LockFreeStack {
	// 调用 newLFS 函数创建一个新的 LockFreeStack 栈，参数为一个新的节点池
	// Call the newLFS function to create a new LockFreeStack stack, the parameter is a new node pool
	return newLFS(shd.NewNodePool())
}

// newLFS 函数用于创建一个新的 LockFreeStack 栈，参数为一个节点池
// The newLFS function is used to create a new LockFreeStack stack, the parameter is a node pool
func newLFS(pool *shd.NodePool) *LockFreeStack {
	// 创建一个新的 Node 结构体实例，值为 nil
	// Create a new Node struct instance, the value is nil
	firstNode := shd.NewNode(nil)

	// 返回一个新的 LockFreeStack 栈，该栈的顶部节点是刚刚创建的节点，节点池为传入的参数
	// Return a new LockFreeStack stack, the top node of this stack is the node just created, and the node pool is the passed in parameter
	return &LockFreeStack{
		pool: pool,
		top:  unsafe.Pointer(firstNode),
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

	// 创建一个新的 Node 结构体实例
	// Create a new Node struct instance
	var node *shd.Node
	if s.pool != nil {
		// 如果节点池不为空，那么从节点池中获取一个节点
		// If the node pool is not nil, then get a node from the node pool
		node = s.pool.Get()
		node.Value = value
	} else {
		// 如果节点池为空，那么创建一个新的节点
		// If the node pool is nil, then create a new node
		node = shd.NewNode(value)
	}

	// 使用无限循环，直到成功推入元素
	// Use an infinite loop until an element is successfully pushed
	for {
		// 获取栈顶元素
		// Get the top element of the stack
		top := shd.LoadNode(&s.top)

		// 设置新节点的下一个元素为当前的栈顶元素
		// Set the next element of the new node to the current top element
		node.Next = unsafe.Pointer(top)

		// 使用 CAS 操作尝试修改栈顶元素
		// Use CAS operation to try to modify the top element
		if shd.CompareAndSwapNode(&s.top, top, node) {
			// 如果成功修改，栈的长度加 1
			// If the modification is successful, the length of the stack is increased by 1
			atomic.AddInt64(&s.length, 1)

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
				// 如果成功修改，栈的长度减 1
				// If the modification is successful, the length of the stack is reduced by 1
				atomic.AddInt64(&s.length, -1)

				// 如果节点池不为空，那么将栈顶元素放回节点池
				// If the node pool is not nil, then put the top element back into the node pool
				if s.pool != nil {
					s.pool.Put(top)
				} else {
					// 如果节点池为空，那么重置栈顶元素
					// If the node pool is nil, then reset the top element
					top.ResetAll()
				}

				// 如果结果不是空值，返回结果
				// If the result is not an empty value, return the result
				return result
			}
		}
	}
}

// Length 方法用于获取 LockFreeQueue 队列的长度
// The Length method is used to get the length of the LockFreeQueue queue
func (s *LockFreeStack) Length() int64 {
	// 使用 atomic.Loadint64 函数获取队列的长度
	// Use the atomic.Loadint64 function to get the length of the queue
	return atomic.LoadInt64(&s.length)
}

// IsEmpty 方法用于判断 LockFreeStack 栈是否为空
// The IsEmpty method is used to determine whether the LockFreeStack stack is empty
func (s *LockFreeStack) IsEmpty() bool {
	// 使用 Length 方法获取栈的长度，如果长度为 0，那么栈为空
	// Use the Length method to get the length of the stack, if the length is 0, then the stack is empty
	return s.Length() == 0
}

// Reset 方法用于重置 LockFreeQueue 队列
// The Reset method is used to reset the LockFreeQueue queue
func (s *LockFreeStack) Reset() {
	// 将队列的头节点和尾节点都设置为新创建的节点
	// Set both the head node and the tail node of the queue to the newly created node
	s.top = unsafe.Pointer(shd.NewNode(nil))

	// 使用 atomic.Storeint64 函数将队列的长度设置为 0
	// Use the atomic.Storeint64 function to set the length of the queue to 0
	atomic.StoreInt64(&s.length, 0)
}
