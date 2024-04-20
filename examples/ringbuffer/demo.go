package main

import (
	"fmt"

	"github.com/shengyanli1982/lockfree/ringbuffer"
)

func main() {
	// 使用 ringbuffer.New 函数创建一个新的环形缓冲区，容量为 10
	// Create a new ring buffer with a capacity of 10 using the ringbuffer.New function
	rb := ringbuffer.New(10)

	// 使用 for 循环向环形缓冲区中推入 10 个元素
	// Use a for loop to push 10 elements into the ring buffer
	for i := 0; i < 10; i++ {
		// 使用 Push 方法尝试推入元素
		// Use the Push method to try to push an element
		if ok := rb.Push(i); !ok {
			// 如果推入失败，打印失败信息
			// If the push fails, print the failure message
			fmt.Printf(">> push failed: %v\n", i)
		}
	}

	// 使用 for 循环从环形缓冲区中弹出 10 个元素
	// Use a for loop to pop 10 elements from the ring buffer
	for i := 0; i < 10; i++ {
		// 使用 Pop 方法尝试弹出元素
		// Use the Pop method to try to pop an element
		if v, ok := rb.Pop(); ok {
			// 如果弹出成功，打印弹出的元素
			// If the pop is successful, print the popped element
			fmt.Printf(">> pop: %v\n", v)
		} else {
			// 如果弹出失败，打印失败信息
			// If the pop fails, print the failure message
			fmt.Printf(">> pop failed: %v\n", i)
		}
	}
}
