package queue

import (
	"sync/atomic"
	"unsafe"

	shd "github.com/shengyanli1982/lockfree/internal/shared"
)

// LockFreeQueue 是一个无锁队列结构体
// LockFreeQueue is a lock-free queue struct
type LockFreeQueue struct {
	// length 是队列的长度
	// length is the length of the queue
	length int64

	// head 是指向队列头部的指针
	// head is a pointer to the head of the queue
	head unsafe.Pointer

	// tail 是指向队列尾部的指针
	// tail is a pointer to the tail of the queue
	tail unsafe.Pointer
}

// New 函数用于创建一个新的 LockFreeQueue 结构体实例
// The New function is used to create a new instance of the LockFreeQueue struct
func New() *LockFreeQueue {
	// 创建一个新的 Node 结构体实例
	// Create a new Node struct instance
	emptyNode := shd.NewNode(shd.EmptyValue)

	// 返回一个新的 LockFreeQueue 结构体实例，其中 head 和 tail 都指向 EmptyNode 节点
	// Returns a new instance of the LockFreeQueue struct, where both head and tail point to the dummy node
	return &LockFreeQueue{
		head: unsafe.Pointer(emptyNode),
		tail: unsafe.Pointer(emptyNode),
	}
}

// Push 方法用于将一个值添加到 LockFreeQueue 队列的末尾
// The Push method is used to add a value to the end of the LockFreeQueue queue
func (q *LockFreeQueue) Push(value interface{}) {
	// 检查值是否为空, 如果为空则直接返回
	// Check if the value is nil, if it is, return directly
	if value == nil {
		return
	}

	// 创建一个新的 Node 结构体实例
	// Create a new Node struct instance
	node := shd.NewNode(value)

	// 使用无限循环来尝试将新节点添加到队列的末尾
	// Use an infinite loop to try to add the new node to the end of the queue
	for {
		// 加载队列的尾节点
		// Load the tail node of the queue
		tail := shd.LoadNode(&q.tail)

		// 加载尾节点的下一个节点
		// Load the next node of the tail node
		next := shd.LoadNode(&tail.Next)

		// 检查尾节点是否仍然是队列的尾节点
		// Check if the tail node is still the tail node of the queue
		if tail == shd.LoadNode(&q.tail) {
			// 如果尾节点的下一个节点是 nil，说明尾节点是队列的最后一个节点
			// If the next node of the tail node is nil, it means that the tail node is the last node of the queue
			if next == nil {
				// 尝试将尾节点的下一个节点设置为新节点
				// Try to set the next node of the tail node to the new node
				if shd.CompareAndSwapNode(&tail.Next, next, node) {
					// 如果成功，那么将队列的尾节点设置为新节点
					// If successful, then set the tail node of the queue to the new node
					shd.CompareAndSwapNode(&q.tail, tail, node)

					// 并增加队列的长度
					// And increase the length of the queue
					atomic.AddInt64(&q.length, 1)

					// 然后返回，结束函数
					// Then return to end the function
					return
				}
			} else {
				// 如果尾节点的下一个节点不是 nil，说明尾节点不是队列的最后一个节点，那么将队列的尾节点设置为尾节点的下一个节点
				// If the next node of the tail node is not nil, it means that the tail node is not the last node of the queue, then set the tail node of the queue to the next node of the tail node
				shd.CompareAndSwapNode(&q.tail, tail, next)
			}
		}
	}
}

// Pop 方法用于从 LockFreeQueue 队列的头部移除并返回一个值
// The Pop method is used to remove and return a value from the head of the LockFreeQueue queue
func (q *LockFreeQueue) Pop() interface{} {
	// 使用无限循环来尝试从队列的头部移除一个值
	// Use an infinite loop to try to remove a value from the head of the queue
	for {
		// 加载队列的头节点
		// Load the head node of the queue
		head := shd.LoadNode(&q.head)

		// 加载队列的尾节点
		// Load the tail node of the queue
		tail := shd.LoadNode(&q.tail)

		// 加载头节点的下一个节点
		// Load the next node of the head node
		next := shd.LoadNode(&head.Next)

		// 检查头节点是否仍然是队列的头节点
		// Check if the head node is still the head node of the queue
		if head == shd.LoadNode(&q.head) {
			// 如果头节点等于尾节点
			// If the head node is equal to the tail node
			if head == tail {
				// 如果头节点的下一个节点是 nil，说明队列是空的，返回 nil
				// If the next node of the head node is nil, it means that the queue is empty, return nil
				if next == nil {
					return nil
				}

				// 如果头节点的下一个节点不是 nil，说明尾节点落后了，尝试将队列的尾节点设置为头节点的下一个节点
				// If the next node of the head node is not nil, it means that the tail node is lagging behind, try to set the tail node of the queue to the next node of the head node
				shd.CompareAndSwapNode(&q.tail, tail, next)
			} else {
				// 并返回头节点的值
				// And return the value of the head node
				result := next.Value

				// 如果头节点不等于尾节点，尝试将队列的头节点设置为头节点的下一个节点
				// If the head node is not equal to the tail node, try to set the head node of the queue to the next node of the head node
				if shd.CompareAndSwapNode(&q.head, head, next) {

					// 如果成功，那么减少队列的长度
					// If successful, then decrease the length of the queue
					atomic.AddInt64(&q.length, -1)

					// 然后重置头节点
					// Then reset the head node
					shd.ResetNodeAll(head)

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
}

// Length 方法用于获取 LockFreeQueue 队列的长度
// The Length method is used to get the length of the LockFreeQueue queue
func (q *LockFreeQueue) Length() int64 {
	// 使用 atomic.Loadint64 函数获取队列的长度
	// Use the atomic.Loadint64 function to get the length of the queue
	return atomic.LoadInt64(&q.length)
}

// Reset 方法用于重置 LockFreeQueue 队列
// The Reset method is used to reset the LockFreeQueue queue
func (q *LockFreeQueue) Reset() {
	// 创建一个新的 Node 结构体实例
	// Create a new Node struct instance
	emptyNode := shd.NewNode(shd.EmptyValue)

	// 将队列的头节点和尾节点都设置为新创建的节点
	// Set both the head node and the tail node of the queue to the newly created node
	q.head = unsafe.Pointer(emptyNode)
	q.tail = unsafe.Pointer(emptyNode)

	// 使用 atomic.Storeint64 函数将队列的长度设置为 0
	// Use the atomic.Storeint64 function to set the length of the queue to 0
	atomic.StoreInt64(&q.length, 0)
}

func (q *LockFreeQueue) IsEmpty() bool {
	// 使用 atomic.LoadInt64 函数获取队列的长度，如果长度为 0，那么队列为空
	// Use the atomic.LoadInt64 function to get the length of the queue, if the length is 0, then the queue is empty
	return atomic.LoadInt64(&q.length) == 0
}
