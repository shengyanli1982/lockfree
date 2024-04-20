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
