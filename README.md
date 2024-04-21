English | [中文](./README_CN.md)

<div align="center">
	<img src="assets/logo.png" alt="logo" width="500px">
</div>

[![Go Report Card](https://goreportcard.com/badge/github.com/shengyanli1982/lockfree)](https://goreportcard.com/report/github.com/shengyanli1982/lockfree)
[![Build Status](https://github.com/shengyanli1982/lockfree/actions/workflows/test.yaml/badge.svg)](https://github.com/shengyanli1982/lockfree/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/shengyanli1982/lockfree.svg)](https://pkg.go.dev/github.com/shengyanli1982/lockfree)

# Introduction

`lockfree` is a Go library that provides lock-free data structures and algorithms. It is designed to be simple, easy to use, and high performance, making it suitable for high-concurrency scenarios.

**Currently, `lockfree` offers the following data structures:**

-   `Queue`: A lock-free queue
-   `Stack`: A lock-free stack
-   `RingBuffer`: A lock-free ring buffer

# Why use `lockfree`?

In high-concurrency scenarios, traditional lock-based data structures can introduce performance bottlenecks due to lock contention. Lock-free data structures can overcome this issue.

I created `lockfree` to provide Go developers with a straightforward, user-friendly, and high-performance lock-free library. Through this project, I aim to enhance my skills in this area and help you overcome common challenges.

# Advantages

-   Simple and user-friendly
-   No external dependencies
-   High performance
-   Thread-safe
-   Supports various data types and structures

# Installation

```bash
go get github.com/shengyanli1982/lockfree
```

# Benchmark

The following benchmark results show the performance of the `lockfree` library compared to the standard `channel` package in Go.

| Benchmark                             | Operations | Time/op     | Bytes/op | Allocs/op   |
| ------------------------------------- | ---------- | ----------- | -------- | ----------- |
| BenchmarkStdChannel-8                 | 15,357,115 | 81.55 ns/op | 0 B/op   | 0 allocs/op |
| BenchmarkStdChannelParallel-8         | 10,393,086 | 116.2 ns/op | 0 B/op   | 0 allocs/op |
| BenchmarkLockFreeQueue-8              | 8,799,632  | 126.2 ns/op | 31 B/op  | 1 allocs/op |
| BenchmarkLockFreeQueueParallel-8      | 6,817,446  | 174.0 ns/op | 24 B/op  | 1 allocs/op |
| BenchmarkLockFreeStack-8              | 9,490,305  | 108.8 ns/op | 31 B/op  | 1 allocs/op |
| BenchmarkLockFreeStackParallel-8      | 8,942,202  | 134.4 ns/op | 24 B/op  | 1 allocs/op |
| BenchmarkLockFreeRingBuffer-8         | 12,610,683 | 114.7 ns/op | 20 B/op  | 2 allocs/op |
| BenchmarkLockFreeRingBufferParallel-8 | 6,104,230  | 199.1 ns/op | 21 B/op  | 2 allocs/op |

### Struct Memory Alignment

**1. Queue**

```bash
Queue alignment:

---- Fields in struct ----
+----+----------------+-----------+-----------+
| ID |   FIELDTYPE    | FIELDNAME | FIELDSIZE |
+----+----------------+-----------+-----------+
| A  | int64          | length    | 8         |
| B  | unsafe.Pointer | head      | 8         |
| C  | unsafe.Pointer | tail      | 8         |
+----+----------------+-----------+-----------+
---- Memory layout ----
|A|A|A|A|A|A|A|A|
|B|B|B|B|B|B|B|B|
|C|C|C|C|C|C|C|C|

total cost: 24 Bytes.
```

**2. Stack**

```bash
Stack alignment:

---- Fields in struct ----
+----+----------------+-----------+-----------+
| ID |   FIELDTYPE    | FIELDNAME | FIELDSIZE |
+----+----------------+-----------+-----------+
| A  | int64          | length    | 8         |
| B  | unsafe.Pointer | top       | 8         |
+----+----------------+-----------+-----------+
---- Memory layout ----
|A|A|A|A|A|A|A|A|
|B|B|B|B|B|B|B|B|

total cost: 16 Bytes.
```

**3. RingBuffer**

```bash
RingBuffer alignment:

---- Fields in struct ----
+----+------------------+-----------+-----------+
| ID |    FIELDTYPE     | FIELDNAME | FIELDSIZE |
+----+------------------+-----------+-----------+
| A  | int64            | capacity  | 8         |
| B  | int64            | head      | 8         |
| C  | int64            | tail      | 8         |
| D  | int64            | count     | 8         |
| E  | []unsafe.Pointer | data      | 24        |
+----+------------------+-----------+-----------+
---- Memory layout ----
|A|A|A|A|A|A|A|A|
|B|B|B|B|B|B|B|B|
|C|C|C|C|C|C|C|C|
|D|D|D|D|D|D|D|D|
|E|E|E|E|E|E|E|E|
|E|E|E|E|E|E|E|E|
|E|E|E|E|E|E|E|E|

total cost: 56 Bytes.
```

# Quick Start

`lockfree` is designed to be easy to use. It provides a simple interface and follows good functional packaging principles, allowing users to quickly get started without requiring extensive learning or training.

## 1. Queue

The `LockFreeQueue` is a thread-safe and lock-free `fifo` data structure. It offers basic operations without support for delaying, priority, timeout, or blocking operations. It is designed to be very simple.

### Methods

-   `Push`: Pushes an element into the queue
-   `Pop`: Pops an element from the queue
-   `Length`: Gets the number of elements in the queue
-   `IsEmpty`: Checks if the queue is empty
-   `Reset`: Resets the queue

### Example

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

**Result**

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

## 2. Stack

The `LockFreeStack` is a thread-safe and lock-free `lifo` data structure. It provides simple methods for pushing and popping elements, as well as getting the length and checking if the stack is empty.

### Methods

-   `Push`: Pushes an element onto the stack
-   `Pop`: Pops an element from the stack
-   `Length`: Gets the number of elements in the stack
-   `IsEmpty`: Checks if the stack is empty
-   `Reset`: Resets the stack

### Example

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

**Result**

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

## 3. RingBuffer

The `LockFreeRingBuffer` is a thread-safe and lock-free data structure that implements a ring buffer. It provides methods for pushing and popping elements, as well as getting the length and checking if the buffer is full or empty.

### Methods

-   `Push`: Pushes an element into the ring buffer
-   `Pop`: Pops an element from the ring buffer
-   `Count`: Gets the number of elements in the ring buffer
-   `Reset`: Resets the ring buffer
-   `IsFull`: Checks if the ring buffer is full
-   `IsEmpty`: Checks if the ring buffer is empty

### Example

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

**Result**

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
