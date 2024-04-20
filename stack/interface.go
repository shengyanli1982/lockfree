package stack

// Stack 是一个接口，定义了堆栈的基本操作
// Stack is an interface that defines basic operations of a stack
type Stack = interface {
	// Push 方法用于向堆栈中添加一个元素
	// The Push method is used to add an element to the stack
	Push(value interface{})

	// Pop 方法用于从堆栈中移除并返回一个元素
	// The Pop method is used to remove and return an element from the stack
	Pop() interface{}

	// Reset 方法用于重置/清空堆栈
	// The Reset method is used to reset/clear the stack
	Reset()

	// Length 方法返回堆栈的长度（即堆栈中元素的数量）
	// The Length method returns the length of the stack (i.e., the number of elements in the stack)
	Length() int64

	// IsEmpty 方法用于检查堆栈是否为空
	// The IsEmpty method is used to check if the stack is empty
	IsEmpty() bool
}
