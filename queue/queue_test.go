package queue

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLockFreeQueue_Standard(t *testing.T) {
	// Number of elements to test
	count := 1000000

	// Create a new queue
	q := New()

	// Test enqueueing elements into the queue
	for i := 0; i < count; i++ {
		q.Push(i)
	}

	// Verify the elements in the queue
	for i := 0; i < count; i++ {
		v := q.Pop()
		if v != nil {
			assert.Equal(t, i, v, "Incorrect value in the queue. Expected %d, got %d", i, v)
		}
	}

	// Verify the queue length
	assert.Equal(t, int64(0), q.Length(), "Incorrect queue length. Expected 0, got %d", q.Length())
}

func TestLockFreeQueue_Length(t *testing.T) {
	q := New()

	// Test the length of an empty queue
	assert.Equal(t, int64(0), q.Length(), "Incorrect queue length. Expected 0, got %d", q.Length())

	// Test the length of a non-empty queue
	for i := 0; i < 100; i++ {
		q.Push(i)
		assert.Equal(t, int64(i+1), q.Length(), "Incorrect queue length. Expected %d, got %d", i+1, q.Length())
	}

	// Test the length of a queue after popping elements
	for i := 0; i < 100; i++ {
		q.Pop()
		assert.Equal(t, int64(100-i-1), q.Length(), "Incorrect queue length. Expected %d, got %d", 100-i-1, q.Length())
	}
}

func TestLockFreeQueue_EmptyPop(t *testing.T) {
	q := New()

	// Test popping elements from an empty queue
	for i := 0; i < 100; i++ {
		v := q.Pop()
		assert.Nil(t, v, "Expected nil value from an empty queue")
	}
}

func TestLockFreeQueue_Parallel(t *testing.T) {
	nums := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	q := New()

	// Test enqueueing elements into the queue
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			q.Push(i)
		}(i)
	}
	wg.Wait()

	// Verify the elements in the queue
	wg = sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			v := q.Pop()
			if v != nil && v.(int) != i {
				assert.Contains(t, nums, v, "Incorrect value in the queue. Expected %d, got %d", i, v)
			}
		}(i)
	}
	wg.Wait()
}

func TestLockFreeQueue_ParallelAtSametime(t *testing.T) {
	nums := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	q := New()

	// Test enqueueing elements into the queue
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			q.Push(i)
		}(i)
	}

	// Verify the elements in the queue
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			v := q.Pop()
			if v != nil && v.(int) != i {
				assert.Contains(t, nums, v.(int), "Incorrect value in the queue. Expected %d, got %d", i, v)
			}
		}(i)
	}
	wg.Wait()
}

func TestLockFreeQueue_ParallelDevilMode(t *testing.T) {
	nums := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	q := New()

	// Test enqueueing elements into the queue
	wg := sync.WaitGroup{}
	for j := 0; j < 10000; j++ {
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				q.Push(i)
			}(i)
		}
	}

	// Verify the elements in the queue
	for j := 0; j < 10000; j++ {
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				v := q.Pop()
				if v != nil && v.(int) != i {
					assert.Contains(t, nums, v.(int), "Incorrect value in the queue. Expected %d, got %d", i, v)
				}
			}(i)
		}
	}
	wg.Wait()
}
