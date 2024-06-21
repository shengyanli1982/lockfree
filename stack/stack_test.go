package stack

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLockFreeStack_Standard(t *testing.T) {
	// Number of elements to test
	count := 1000000

	// Create a new stack
	q := New()

	// Test enstacking elements into the stack
	for i := 0; i < count; i++ {
		q.Push(i)
	}

	// Verify the stack length
	assert.Equal(t, int64(count), q.Length(), "Incorrect stack length. Expected %d, got %d", count, q.Length())

	// Verify the elements in the stack
	for i := count - 1; i >= 0; i-- {
		v := q.Pop()
		if v != nil {
			assert.Equal(t, i, v, "Incorrect value in the stack. Expected %d, got %d", i, v)
		}
	}

	// Verify the stack length
	assert.Equal(t, int64(0), q.Length(), "Incorrect stack length. Expected 0, got %d", q.Length())
}

func TestLockFreeStack_Length(t *testing.T) {
	q := New()

	// Test the length of an empty stack
	assert.Equal(t, int64(0), q.Length(), "Incorrect stack length. Expected 0, got %d", q.Length())

	// Test the length of a non-empty stack
	for i := 0; i < 100; i++ {
		q.Push(i)
		assert.Equal(t, int64(i+1), q.Length(), "Incorrect stack length. Expected %d, got %d", i+1, q.Length())
	}

	// Test the length of a stack after popping elements
	for i := 0; i < 100; i++ {
		q.Pop()
		assert.Equal(t, int64(100-i-1), q.Length(), "Incorrect stack length. Expected %d, got %d", 100-i-1, q.Length())
	}
}

func TestLockFreeStack_EmptyPop(t *testing.T) {
	q := New()

	// Test popping elements from an empty stack
	for i := 0; i < 100; i++ {
		v := q.Pop()
		assert.Nil(t, v, "Expected nil value from an empty stack")
	}
}

func TestLockFreeStack_Parallel(t *testing.T) {
	nums := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	q := New()

	// Test enstacking elements into the stack
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			q.Push(i)
		}(i)
	}
	wg.Wait()

	// Verify the elements in the stack
	wg = sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			v := q.Pop()
			if v != nil && v.(int) != i {
				assert.Contains(t, nums, v, "Incorrect value in the stack. Expected %d, got %d", i, v)
			}
		}(i)
	}
	wg.Wait()
}

func TestLockFreeStack_ParallelAtSametime(t *testing.T) {
	nums := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	q := New()

	// Test enstacking elements into the stack
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			q.Push(i)
		}(i)
	}

	// Verify the elements in the stack
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			v := q.Pop()
			if v != nil && v.(int) != i {
				assert.Contains(t, nums, v.(int), "Incorrect value in the stack. Expected %d, got %d", i, v)
			}
		}(i)
	}
	wg.Wait()
}

func TestLockFreeStack_ParallelDevilMode(t *testing.T) {
	nums := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	q := New()

	// Test enstacking elements into the stack
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

	// Verify the elements in the stack
	for j := 0; j < 10000; j++ {
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				v := q.Pop()
				if v != nil && v.(int) != i {
					assert.Contains(t, nums, v.(int), "Incorrect value in the stack. Expected %d, got %d", i, v)
				}
			}(i)
		}
	}
	wg.Wait()
}

func TestLockFreeStack_WithPool_Standard(t *testing.T) {
	// Number of elements to test
	count := 1000000

	// Create a new stack
	q := NewWithPool()

	// Test enstacking elements into the stack
	for i := 0; i < count; i++ {
		q.Push(i)
	}

	// Verify the stack length
	assert.Equal(t, int64(count), q.Length(), "Incorrect stack length. Expected %d, got %d", count, q.Length())

	// Verify the elements in the stack
	for i := count - 1; i >= 0; i-- {
		v := q.Pop()
		if v != nil {
			assert.Equal(t, i, v, "Incorrect value in the stack. Expected %d, got %d", i, v)
		}
	}

	// Verify the stack length
	assert.Equal(t, int64(0), q.Length(), "Incorrect stack length. Expected 0, got %d", q.Length())
}

func TestLockFreeStack_WithPool_Length(t *testing.T) {
	q := NewWithPool()

	// Test the length of an empty stack
	assert.Equal(t, int64(0), q.Length(), "Incorrect stack length. Expected 0, got %d", q.Length())

	// Test the length of a non-empty stack
	for i := 0; i < 100; i++ {
		q.Push(i)
		assert.Equal(t, int64(i+1), q.Length(), "Incorrect stack length. Expected %d, got %d", i+1, q.Length())
	}

	// Test the length of a stack after popping elements
	for i := 0; i < 100; i++ {
		q.Pop()
		assert.Equal(t, int64(100-i-1), q.Length(), "Incorrect stack length. Expected %d, got %d", 100-i-1, q.Length())
	}
}

func TestLockFreeStack_WithPool_EmptyPop(t *testing.T) {
	q := NewWithPool()

	// Test popping elements from an empty stack
	for i := 0; i < 100; i++ {
		v := q.Pop()
		assert.Nil(t, v, "Expected nil value from an empty stack")
	}
}

func TestLockFreeStack_WithPool_Parallel(t *testing.T) {
	nums := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	q := NewWithPool()

	// Test enstacking elements into the stack
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			q.Push(i)
		}(i)
	}
	wg.Wait()

	// Verify the elements in the stack
	wg = sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			v := q.Pop()
			if v != nil && v.(int) != i {
				assert.Contains(t, nums, v, "Incorrect value in the stack. Expected %d, got %d", i, v)
			}
		}(i)
	}
	wg.Wait()
}

func TestLockFreeStack_WithPool_ParallelAtSametime(t *testing.T) {
	nums := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	q := NewWithPool()

	// Test enstacking elements into the stack
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			q.Push(i)
		}(i)
	}

	// Verify the elements in the stack
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			v := q.Pop()
			if v != nil && v.(int) != i {
				assert.Contains(t, nums, v.(int), "Incorrect value in the stack. Expected %d, got %d", i, v)
			}
		}(i)
	}
	wg.Wait()
}

func TestLockFreeStack_WithPool_ParallelDevilMode(t *testing.T) {
	nums := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	q := NewWithPool()

	// Test enstacking elements into the stack
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

	// Verify the elements in the stack
	for j := 0; j < 10000; j++ {
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				v := q.Pop()
				if v != nil && v.(int) != i {
					assert.Contains(t, nums, v.(int), "Incorrect value in the stack. Expected %d, got %d", i, v)
				}
			}(i)
		}
	}
	wg.Wait()
}
