package main

import (
	"fmt"

	shd "github.com/shengyanli1982/lockfree/internal/shared"
	memalign "github.com/vearne/mem-align"
)

func main() {
	fmt.Printf("Node struct alignment:\n\n")
	memalign.PrintStructAlignment(shd.Node{})
}
