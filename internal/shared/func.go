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

// SetNode 函数用于设置一个节点的指针
// The SetNode function is used to set a pointer to a node
func SetNode(p *unsafe.Pointer, node *Node) {
	// 使用 atomic.StorePointer 方法存储一个指向节点的指针
	// Use the atomic.StorePointer method to store a pointer to a node
	atomic.StorePointer(p, unsafe.Pointer(node))
}

// CompareAndSwapInt64 函数用于比较并交换一个 int64 类型的值
// The CompareAndSwapInt64 function is used to compare and swap an int64 type value
func CompareAndSwapInt64(p *int64, old, new int64) bool {
	// 使用 atomic.CompareAndSwapInt64 方法比较并交换一个 int64 类型的值
	// Use the atomic.CompareAndSwapInt64 method to compare and swap an int64 type value
	return atomic.CompareAndSwapInt64(p, old, new)
}
