package shared

import (
	"unsafe"
)

// Node 数据单元节点
// Node represents a data unit node
type Node struct {
	// Value 是节点存储的值，类型为 interface{}，可以存储任何类型的值
	// Value is the Value stored in the node, of type interface{}, which can store any type of Value
	Value interface{}

	// Index 是节点的索引，类型为 uint64
	// Index is the index of the node, of type uint64
	Index uint64

	// Next 是指向下一个节点的指针，类型为 unsafe.Pointer
	// Next is a pointer to the Next node, of type unsafe.Pointer
	Next unsafe.Pointer
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

	// 将 index 字段设置为 0
	// Set the index field to 0
	n.Index = 0

	// 将 next 字段设置为 nil
	// Set the next field to nil
	n.Next = nil
}

func (n *Node) ResetValue() {
	// 将 value 字段设置为 nil
	// Set the value field to nil
	n.Value = nil
}

// 定义一个空的结构体，没有任何字段
// Define an empty struct, with no fields
var EmptyValue = struct{}{}

// DefaultRingSize 是环形队列的默认大小
// DefaultRingSize is the default size of the ring queue
const DefaultRingSize = 2048
