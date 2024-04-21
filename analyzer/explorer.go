package main

import (
	"fmt"

	shd "github.com/shengyanli1982/lockfree/internal/shared"
	"github.com/shengyanli1982/lockfree/queue"
	"github.com/shengyanli1982/lockfree/ringbuffer"
	"github.com/shengyanli1982/lockfree/stack"
	memalign "github.com/vearne/mem-align"
)

func main() {
	fmt.Printf("Node alignment:\n\n")
	memalign.PrintStructAlignment(shd.Node{})
	fmt.Printf("\n =========================== \n\n")

	fmt.Printf("Queue alignment:\n\n")
	memalign.PrintStructAlignment(queue.LockFreeQueue{})
	fmt.Printf("\n =========================== \n\n")

	fmt.Printf("Stack alignment:\n\n")
	memalign.PrintStructAlignment(stack.LockFreeStack{})
	fmt.Printf("\n =========================== \n\n")

	fmt.Printf("RingBuffer alignment:\n\n")
	memalign.PrintStructAlignment(ringbuffer.LockFreeRingBuffer{})
}
