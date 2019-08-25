package semaphore // import "github.com/Warashi/go-semaphore"

import (
	"sync"
	"sync/atomic"
)

type Semaphore interface {
	Acquire()
	Release()
	IncreaseSize(n int64) int64
}

var _ Semaphore = &NonFairSemaphore{}

type NonFairSemaphore struct {
	size int64
	cur  int64
	cond *sync.Cond
}

func NewNonFair(size int64) *NonFairSemaphore {
	return &NonFairSemaphore{
		size: size,
		cur:  0,
		cond: sync.NewCond(new(sync.Mutex)),
	}
}

func (s *NonFairSemaphore) Acquire() {
	for {
		if atomic.AddInt64(&s.cur, 1) <= atomic.LoadInt64(&s.size) {
			return
		}
		s.cond.L.Lock()
		atomic.AddInt64(&s.cur, -1)
		for atomic.LoadInt64(&s.size) <= atomic.LoadInt64(&s.cur) {
			s.cond.Wait()
		}
		s.cond.L.Unlock()
	}
}

func (s *NonFairSemaphore) Release() {
	s.cond.L.Lock()
	if atomic.AddInt64(&s.cur, -1) < 0 {
		panic("released semaphore more than acquired")
	}
	s.cond.Signal()
	s.cond.L.Unlock()
}

func (s *NonFairSemaphore) IncreaseSize(n int64) int64 {
	s.cond.L.Lock()
	newSize := atomic.AddInt64(&s.size, n)
	if 1 <= newSize {
		s.cond.Broadcast()
		s.cond.L.Unlock()
		return newSize
	}
	atomic.StoreInt64(&s.size, 1)
	s.cond.Broadcast()
	s.cond.L.Unlock()
	return 1
}
