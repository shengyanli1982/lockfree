package main

import (
	"fmt"

	"github.com/shengyanli1982/lockfree/stack"
)

func main() {
	// 使用 stack.New 函数创建一个新的栈
	// Create a new stack using the stack.New function
	s := stack.New()

	// 使用 for 循环向栈中推入 10 个元素
	// Use a for loop to push 10 elements into the stack
	for i := 0; i < 10; i++ {
		// 使用 Push 方法推入元素
		// Use the Push method to push an element
		s.Push(i)
	}

	// 使用 for 循环从栈中弹出 10 个元素
	// Use a for loop to pop 10 elements from the stack
	for i := 0; i < 10; i++ {
		// 使用 Pop 方法尝试弹出元素
		// Use the Pop method to try to pop an element
		if v := s.Pop(); v != nil {
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
