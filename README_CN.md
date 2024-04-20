[English](./README.md) | 中文

<div align="center">
	<img src="assets/logo.png" alt="logo" width="500px">
</div>

[![Go Report Card](https://goreportcard.com/badge/github.com/shengyanli1982/lockfree)](https://goreportcard.com/report/github.com/shengyanli1982/lockfree)
[![Build Status](https://github.com/shengyanli1982/lockfree/actions/workflows/test.yaml/badge.svg)](https://github.com/shengyanli1982/lockfree/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/shengyanli1982/lockfree.svg)](https://pkg.go.dev/github.com/shengyanli1982/lockfree)

# 介绍

`lockfree` 是一个提供无锁数据结构和算法的 Go 库。它被设计为简单、易于使用和高性能，适用于高并发场景。

**目前，`lockfree` 提供以下数据结构：**

-   `Queue`：无锁队列
-   `Stack`：无锁栈
-   `RingBuffer`：无锁环形缓冲区

# 为什么使用 `lockfree`？

在高并发场景中，传统的基于锁的数据结构可能由于锁竞争而引入性能瓶颈。无锁数据结构可以克服这个问题。

我创建了 `lockfree` 来为 Go 开发者提供一个简单、用户友好且高性能的无锁库。通过这个项目，我旨在提升自己在这个领域的技能，并帮助您克服常见的挑战。

# 优势

-   简单且用户友好
-   无外部依赖
-   高性能
-   线程安全
-   支持各种数据类型和结构

# 安装

```bash
go get github.com/shengyanli1982/lockfree
```

# 基准测试

以下基准测试结果展示了 `lockfree` 库与 Go 中标准的 `channel` 包相比的性能表现。

```bash
$ go test -bench=.
goos: darwin
goarch: amd64
pkg: github.com/shengyanli1982/lockfree/benchmark
cpu: Intel(R) Xeon(R) CPU E5-4627 v2 @ 3.30GHz
BenchmarkStdChannel-8                   	13972540	        77.07 ns/op	       0 B/op	       0 allocs/op
BenchmarkStdChannelParallel-8           	10272735	       113.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkLockFreeQueue-8                	 9709624	       126.0 ns/op	      39 B/op	       1 allocs/op
BenchmarkLockFreeQueueParallel-8        	 4724350	       253.9 ns/op	      32 B/op	       1 allocs/op
BenchmarkLockFreeStack-8                	10666360	       107.7 ns/op	      39 B/op	       1 allocs/op
BenchmarkLockFreeStackParallel-8        	 4609512	       256.9 ns/op	      32 B/op	       1 allocs/op
BenchmarkLockFreeRingBuffer-8           	10800597	       111.1 ns/op	      24 B/op	       3 allocs/op
BenchmarkLockFreeRingBufferParallel-8   	 4838744	       249.5 ns/op	      48 B/op	       6 allocs/op
```

# 快速入门

`lockfree` 的设计目标是易于使用。它提供了简单的接口，并遵循良好的功能封装原则，使用户能够快速入门，无需进行大量的学习或培训。

## 1. 队列

`LockFreeQueue` 是一个线程安全且无锁的 `fifo` 数据结构。它提供了基本的操作，但不支持延迟、优先级、超时或阻塞操作。它的设计非常简单。

### 方法

-   `Push`：将元素推入队列
-   `Pop`：从队列中弹出元素
-   `Length`：获取队列中的元素数量
-   `IsEmpty`：检查队列是否为空
-   `Reset`：重置队列

### 示例

```go
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
```

**执行结果**

```bash
$ go run demo.go
>> pop: 0
>> pop: 1
>> pop: 2
>> pop: 3
>> pop: 4
>> pop: 5
>> pop: 6
>> pop: 7
>> pop: 8
>> pop: 9
```

## 2. 栈

`LockFreeStack` 是一个线程安全且无锁的 `lifo` 数据结构。它提供了简单的方法来推入和弹出元素，以及获取栈的长度和检查栈是否为空。

### 方法

-   `Push`：将元素推入栈
-   `Pop`：从栈中弹出元素
-   `Length`：获取栈中的元素数量
-   `IsEmpty`：检查栈是否为空
-   `Reset`：重置栈

### 示例

```go
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
```

**执行结果**

```bash
$ go run demo.go
>> pop: 9
>> pop: 8
>> pop: 7
>> pop: 6
>> pop: 5
>> pop: 4
>> pop: 3
>> pop: 2
>> pop: 1
>> pop: 0
```

## 3. 环形缓冲区

`LockFreeRingBuffer` 是一个线程安全且无锁的数据结构，实现了环形缓冲区。它提供了推入和弹出元素的方法，以及获取缓冲区长度和检查缓冲区是否满或空的功能。

### 方法

-   `Push`：将元素推入环形缓冲区
-   `Pop`：从环形缓冲区弹出元素
-   `Count`：获取环形缓冲区中的元素数量
-   `Reset`：重置环形缓冲区
-   `IsFull`：检查环形缓冲区是否已满
-   `IsEmpty`：检查环形缓冲区是否为空

### 示例

```go
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
```

**执行结果**

```bash
$ go run demo.go
>> pop: 0
>> pop: 1
>> pop: 2
>> pop: 3
>> pop: 4
>> pop: 5
>> pop: 6
>> pop: 7
>> pop: 8
>> pop: 9
```
