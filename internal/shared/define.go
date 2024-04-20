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

	// Index 是节点的索引，类型为 int64
	// Index is the index of the node, of type int64
	Index int64

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

// 定义一个空的结构体，没有任何字段
// Define an empty struct, with no fields
var EmptyValue = struct{}{}
