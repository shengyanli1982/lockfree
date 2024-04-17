package stack

// Stack 是一个接口，定义了队列的基本操作
// Stack is an interface that defines basic operations of a queue
type Stack = interface {
	// Push 方法用于向队列中添加一个元素
	// The Push method is used to add an element to the queue
	Push(value interface{})

	// Pop 方法用于从队列中移除并返回一个元素
	// The Pop method is used to remove and return an element from the queue
	Pop() interface{}

	// Reset 方法用于重置/清空队列
	// The Reset method is used to reset/clear the queue
	Reset()

	// Length 方法返回队列的长度（即队列中元素的数量）
	// The Length method returns the length of the queue (i.e., the number of elements in the queue)
	Length() uint64
}
