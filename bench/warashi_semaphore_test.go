package bench

import (
	"testing"

	"github.com/Warashi/go-semaphore"
)

func BenchmarkNonFairSemaphore_Acquire(b *testing.B) {
	const size = 1000
	s := semaphore.NewNonFair(size)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.Acquire()
			s.Release()
		}
	})
}
