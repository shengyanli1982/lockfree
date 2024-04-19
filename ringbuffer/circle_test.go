package ringbuffer

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLockFreeRingBuffer_Push(t *testing.T) {
	r := New(5) // Replace with your desired capacity

	// Push values into the ring buffer
	for i := 0; i < 5; i++ {
		if !r.Push(i) {
			assert.Fail(t, "Failed to push value: %d", i)
		}
	}

	// Try pushing one more value, should fail
	if r.Push(5) {
		assert.Fail(t, "Pushed value when the ring buffer is full")
	}
}

func TestLockFreeRingBuffer_Pop(t *testing.T) {
	r := New(5) // Replace with your desired capacity

	// Push values into the ring buffer
	for i := 0; i < 5; i++ {
		r.Push(i)
	}

	// Pop values from the ring buffer
	for i := 0; i < 5; i++ {
		value, ok := r.Pop()
		assert.True(t, ok, "Failed to pop value")
		assert.Equal(t, i, value, "Incorrect value popped")
	}

	// Try popping one more value, should fail
	_, ok := r.Pop()
	assert.False(t, ok, "Popped value when the ring buffer is empty")
}

func TestLockFreeRingBuffer_Count(t *testing.T) {
	r := New(5) // Replace with your desired capacity

	// Push values into the ring buffer
	for i := 0; i < 5; i++ {
		r.Push(i)
	}

	// Check the count of the ring buffer
	count := r.Count()
	assert.Equal(t, int64(5), count, "Incorrect count of the ring buffer")

	// Pop values from the ring buffer
	for i := 0; i < 5; i++ {
		r.Pop()
	}

	// Check the count of the ring buffer after popping all values
	count = r.Count()
	assert.Equal(t, int64(0), count, "Incorrect count of the ring buffer after popping all values")
}

func TestLockFreeRingBuffer_Standard(t *testing.T) {
	// Test the ring buffer with a large number of elements
	count := 1000000

	r := New(count) // Replace with your desired capacity

	// Test pushing values into the ring buffer
	for i := 0; i < count; i++ {
		if !r.Push(i) {
			assert.Fail(t, "Failed to push value: %d", i)
		}
	}

	// Verify the ring buffer length
	assert.Equal(t, int64(count), r.Count(), "Incorrect ring buffer length. Expected %d, got %d", count, r.Count())

	// Verify the elements in the ring buffer
	for i := 0; i < count; i++ {
		value, ok := r.Pop()
		assert.True(t, ok, "Failed to pop value")
		assert.Equal(t, i, value, "Incorrect value in the ring buffer. Expected %d, got %d", i, value)
	}

	// Verify the ring buffer length
	assert.Equal(t, int64(0), r.Count(), "Incorrect ring buffer length. Expected 0, got %d", r.Count())
}

func TestLockFreeRingBuffer_Length(t *testing.T) {
	r := New(5) // Replace with your desired capacity

	// Test the length of an empty ring buffer
	assert.Equal(t, int64(0), r.Count(), "Incorrect ring buffer length. Expected 0, got %d", r.Count())

	// Test the length of a non-empty ring buffer
	for i := 0; i < 5; i++ {
		r.Push(i)
		assert.Equal(t, int64(i+1), r.Count(), "Incorrect ring buffer length. Expected %d, got %d", i+1, r.Count())
	}

	// Test the length of a ring buffer after popping elements
	for i := 0; i < 5; i++ {
		r.Pop()
		assert.Equal(t, int64(5-i-1), r.Count(), "Incorrect ring buffer length. Expected %d, got %d", 5-i-1, r.Count())
	}

}

func TestLockFreeRingBuffer_EmptyPop(t *testing.T) {
	r := New(5) // Replace with your desired capacity

	// Test popping values from an empty ring buffer
	for i := 0; i < 5; i++ {
		_, ok := r.Pop()
		assert.False(t, ok, "Popped value from an empty ring buffer")
	}
}

func TestLockFreeRingBuffer_Parallel(t *testing.T) {
	nums := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	count := len(nums)
	r := New(count) // Replace with your desired capacity

	// Test enring buffering elements into the ring buffer
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if !r.Push(i) {
				assert.Fail(t, "Failed to push value: %d", i)
			}
		}(i)
	}
	wg.Wait()

	// Verify the elements in the ring buffer
	wg = sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			v, ok := r.Pop()
			assert.True(t, ok, "Failed to pop value")
			if v != nil && v.(int) != i {
				assert.Contains(t, nums, v, "Incorrect value in the ring buffer. Expected %d, got %d", i, v)
			}
		}(i)
	}
	wg.Wait()
}

func TestLockFreeRingBuffer_ParallelAtSametime(t *testing.T) {
	nums := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	count := len(nums)
	r := New(count) // Replace with your desired capacity

	// Test enring buffering elements into the ring buffer
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if !r.Push(i) {
				assert.Fail(t, "Failed to push value: %d", i)
			}
		}(i)
	}

	// Verify the elements in the ring buffer
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			v, ok := r.Pop()
			assert.True(t, ok, "Failed to pop value")
			if v != nil && v.(int) != i {
				assert.Contains(t, nums, v, "Incorrect value in the ring buffer. Expected %d, got %d", i, v)
			}
		}(i)
	}
	wg.Wait()
}

func BenchmarkLockFreeRingBuffer(b *testing.B) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	r := New(b.N)
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
	r := New(5)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Push(1)
			r.Pop()
		}
	})
}

func Benchmark_MOD(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = i % 265
	}
}

func Benchmark_BitOR(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = i & 127
	}
}