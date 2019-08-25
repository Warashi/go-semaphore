package semaphore_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/Warashi/go-semaphore"
)

func TestSemaphore_Release(t *testing.T) {
	s := semaphore.NewNonFair(1)

	var wg sync.WaitGroup
	var c int64
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer s.Release()
			defer wg.Done()
			s.Acquire()
			atomic.AddInt64(&c, 1)
		}()
	}
	wg.Wait()
	if got := atomic.LoadInt64(&c); got != 10 {
		t.Errorf("expected %d, got %d", 10, got)
	}
}

func TestSemaphore_IncreaseSize(t *testing.T) {
	s := semaphore.NewNonFair(0)

	var wg sync.WaitGroup
	var c int64
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Acquire()
			atomic.AddInt64(&c, 1)
		}()
	}
	for i := 0; i < 10; i++ {
		s.IncreaseSize(1)
	}
	wg.Wait()
	if got := atomic.LoadInt64(&c); got != 10 {
		t.Errorf("expected %d, got %d", 10, got)
	}
}

func TestSemaphore_IncreaseSizeBulk(t *testing.T) {
	s := semaphore.NewNonFair(0)

	var wg sync.WaitGroup
	var c int64
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Acquire()
			atomic.AddInt64(&c, 1)
		}()
	}
	s.IncreaseSize(10)
	wg.Wait()
	if got := atomic.LoadInt64(&c); got != 10 {
		t.Errorf("expected %d, got %d", 10, got)
	}
}

func BenchmarkNonFairSemaphore_Acquire(b *testing.B) {
	const size = 10
	s := semaphore.NewNonFair(size)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.Acquire()
			s.Release()
		}
	})
}

func BenchmarkNonFairSemaphore_IncreaseSize(b *testing.B) {
	const size = 10
	s := semaphore.NewNonFair(size)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.Acquire()
			s.IncreaseSize(1)
		}
	})
}
