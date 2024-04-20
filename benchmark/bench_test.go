package benchmark

import (
	"sync"
	"testing"

	"github.com/shengyanli1982/lockfree/queue"
	"github.com/shengyanli1982/lockfree/ringbuffer"
	"github.com/shengyanli1982/lockfree/stack"
)

func BenchmarkStdChannel(b *testing.B) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	ch := make(chan int, b.N)
	b.ResetTimer()
	go func() {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			ch <- i
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			<-ch
		}
	}()
	wg.Wait()
}

func BenchmarkStdChannelParallel(b *testing.B) {
	ch := make(chan int, 5)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ch <- 1
			<-ch
		}
	})
}

func BenchmarkLockFreeQueue(b *testing.B) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	q := queue.New()
	b.ResetTimer()
	go func() {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			q.Push(i)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			q.Pop()
		}
	}()
	wg.Wait()
}

func BenchmarkLockFreeQueueParallel(b *testing.B) {
	q := queue.New()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			q.Push(1)
			q.Pop()
		}
	})
}

func BenchmarkLockFreeStack(b *testing.B) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	q := stack.New()
	b.ResetTimer()
	go func() {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			q.Push(i)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			q.Pop()
		}
	}()
	wg.Wait()
}

func BenchmarkLockFreeStackParallel(b *testing.B) {
	q := stack.New()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			q.Push(1)
			q.Pop()
		}
	})
}

func BenchmarkLockFreeRingBuffer(b *testing.B) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	r := ringbuffer.New(b.N)
	b.ResetTimer()
	go func() {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			r.Push(i)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			r.Pop()
		}
	}()
	wg.Wait()
}

func BenchmarkLockFreeRingBufferParallel(b *testing.B) {
	r := ringbuffer.New(5)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Push(1)
			r.Pop()
		}
	})
}
