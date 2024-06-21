package shared

import (
	"unsafe"
)

// On a 64-bit computer, the size of a Go object is typically a multiple of 8 bytes, which is the size of a pointer on a 64-bit architecture.
// Therefore, the size of a Go object is usually 8 bytes, 16 bytes, 32 bytes, and so on.  24 bytes is not a common size for a Go object.

// Node 数据单元节点
// Node represents a data unit node
type Node struct {
	// Value 是节点存储的值，类型为 interface{}，可以存储任何类型的值
	// Value is the Value stored in the node, of type interface{}, which can store any type of Value
	Value interface{}

	// Next 是指向下一个节点的指针，类型为 unsafe.Pointer
	// Next is a pointer to the Next node, of type unsafe.Pointer
	Next unsafe.Pointer

	// _ 是一个占位符，用于填充内存对齐
	// _ is a placeholder used to fill memory alignment
	_ int64
}

// NewNode 函数用于创建一个新的 Node 结构体实例
// The NewNode function is used to create a new instance of the Node struct
func NewNode(v interface{}) *Node {
	// 返回一个新的 Node 结构体实例
	// Returns a new instance of the Node struct
	return &Node{Value: v}
}

// Reset 方法用于重置 Node 结构体的值
// The Reset method is used to reset the value of the Node struct
func (n *Node) ResetAll() {
	// 将 value 字段设置为 nil
	// Set the value field to nil
	n.Value = nil

	// 将 next 字段设置为 nil
	// Set the next field to nil
	n.Next = nil
}
