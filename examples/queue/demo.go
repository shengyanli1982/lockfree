package main

import (
	"fmt"

	"github.com/shengyanli1982/lockfree/queue"
)

func main() {
	// 使用 queue.New 函数创建一个新的队列
	// Create a new queue using the queue.New function
	q := queue.New()

	// 使用 for 循环向队列中推入 10 个元素
	// Use a for loop to push 10 elements into the queue
	for i := 0; i < 10; i++ {
		// 使用 Push 方法推入元素
		// Use the Push method to push an element
		q.Push(i)
	}

	// 使用 for 循环从队列中弹出 10 个元素
	// Use a for loop to pop 10 elements from the queue
	for i := 0; i < 10; i++ {
		// 使用 Pop 方法尝试弹出元素
		// Use the Pop method to try to pop an element
		if v := q.Pop(); v != nil {
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
