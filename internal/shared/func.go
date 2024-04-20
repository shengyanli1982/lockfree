package shared

import (
	"sync/atomic"
	"unsafe"
)

// LoadNode 函数用于加载指定指针 p 指向的 Node 结构体
// The LoadNode function is used to load the Node struct pointed to by the specified pointer p
func LoadNode(p *unsafe.Pointer) *Node {
	// 使用 atomic.LoadPointer 加载并返回指定指针 p 指向的 Node 结构体
	// Uses atomic.LoadPointer to load and return the Node struct pointed to by the specified pointer p
	return (*Node)(atomic.LoadPointer(p))
}

// CompareAndSwapNode 函数用于比较并交换指定指针 p 指向的 Node 结构体
// The CompareAndSwapNode function is used to compare and swap the Node struct pointed to by the specified pointer p
func CompareAndSwapNode(p *unsafe.Pointer, old, new *Node) bool {
	// 使用 atomic.CompareAndSwapPointer 比较并交换指定指针 p 指向的 Node 结构体
	// Uses atomic.CompareAndSwapPointer to compare and swap the Node struct pointed to by the specified pointer p
	return atomic.CompareAndSwapPointer(p, unsafe.Pointer(old), unsafe.Pointer(new))
}

// Reset 方法用于重置 Node 结构体的值
// The Reset method is used to reset the value of the Node struct
func ResetNodeAll(n *Node) {
	// 将 value 字段设置为 nil
	// Set the value field to nil
	n.Value = nil

	// 将 next 字段设置为 nil
	// Set the next field to nil
	n.Next = nil
}

// ResetNodeValue 方法用于重置 Node 结构体 value 字段的值
// The ResetNodeValue method is used to reset the value of the Node struct value field
func ResetNodeValue(n *Node) {
	// 将 value 字段设置为 nil
	// Set the value field to nil
	n.Value = nil
}

// ResetNodeNext 方法用于重置 Node 结构体 next 字段的值
// The ResetNodeNext method is used to reset the value of the Node struct next field
func ResetNodeNext(n *Node) {
	// 将 next 字段设置为 nil
	// Set the next field to nil
	n.Next = nil
}

// AbsInt64 是一个函数，它接受一个 int64 类型的参数 x，并返回 x 的绝对值
// AbsInt64 is a function that takes an int64 type parameter x and returns the absolute value of x
func AbsInt64(x int64) int64 {
	// 如果 x 小于 0
	// If x is less than 0
	if x < 0 {
		// 使用位运算符 ^ 对 x 进行按位取反，然后加 1，得到 x 的绝对值
		// Use the bitwise operator ^ to bitwise negate x, then add 1 to get the absolute value of x
		return ^x + 1
	}

	// 如果 x 大于或等于 0，直接返回 x
	// If x is greater than or equal to 0, return x directly
	return x
}
