package shared

import (
	"sync"
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

// NodePool 结构体用于表示一个节点池
// The NodePool struct is used to represent a node pool
type NodePool struct {
	// pool 是一个同步池，用于存储和获取节点
	// pool is a sync pool, used to store and retrieve nodes
	pool *sync.Pool
}

// NewNodePool 函数用于创建一个新的节点池
// The NewNodePool function is used to create a new node pool
func NewNodePool() *NodePool {
	return &NodePool{
		// 创建一个新的同步池，当池中没有可用的节点时，会调用 New 方法创建一个新的节点
		// Create a new sync pool, when there are no available nodes in the pool, the New method will be called to create a new node
		pool: &sync.Pool{
			New: func() interface{} {
				return NewNode(nil)
			},
		},
	}
}

// Get 方法用于从节点池中获取一个节点
// The Get method is used to get a node from the node pool
func (np *NodePool) Get() *Node {
	// 使用 sync.Pool 的 Get 方法获取一个节点，然后将其转换为 *Node 类型
	// Use the Get method of sync.Pool to get a node, and then convert it to *Node type
	return np.pool.Get().(*Node)
}

// Put 方法用于将一个节点放回节点池
// The Put method is used to put a node back into the node pool
func (np *NodePool) Put(n *Node) {
	// 如果节点不为 nil
	// If the node is not nil
	if n != nil {
		// 重置节点，包括其值、下一个节点和索引
		// Reset the node, including its value, next node, and index
		n.ResetAll()

		// 使用 sync.Pool 的 Put 方法将节点放回池中
		// Use the Put method of sync.Pool to put the node back into the pool
		np.pool.Put(n)
	}
}
