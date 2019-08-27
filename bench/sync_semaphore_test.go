package bench

import (
	"context"
	"testing"

	"golang.org/x/sync/semaphore"
)

func BenchmarkSyncSemaphore(b *testing.B) {
	const size = 1000
	s := semaphore.NewWeighted(size)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.Acquire(context.Background(), 1)
			s.Release(1)
		}
	})
}
