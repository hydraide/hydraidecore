package lock

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestLockUnlock(t *testing.T) {
	l := New()

	ctx := context.Background()

	lockID, err := l.Lock(ctx, "key1", 2*time.Second)
	assert.NoError(t, err)
	assert.NotEmpty(t, lockID)

	err = l.Unlock("key1", lockID)
	assert.NoError(t, err)
}

func TestParallelLockUnlock(t *testing.T) {
	l := New()

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {

				func() {

					key := fmt.Sprintf("key%d", j)

					ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
					defer cancel()

					lockID, err := l.Lock(ctx, key, 2*time.Second)
					if err != nil {
						t.Errorf("failed to acquire lock: %v", err)
						return
					}
					assert.NotEmpty(t, lockID)

					// Simulate some work with the lock held
					time.Sleep(10 * time.Millisecond)

					err = l.Unlock(key, lockID)
					if err != nil {
						t.Errorf("failed to release lock: %v", err)
						return
					}

				}()

			}
		}()
	}

	wg.Wait()
}

func TestLockTimeout(t *testing.T) {
	l := New()

	ctx := context.Background()

	// Lock megszerzése, ami 1 másodperc múlva lejár
	lockID, err := l.Lock(ctx, "key1", 1*time.Second)
	assert.NoError(t, err)
	assert.NotEmpty(t, lockID)

	time.Sleep(2 * time.Second) // waiting for the lock to expire

	// Próbáljuk meg újra megszerezni a lockot
	lockID2, err := l.Lock(ctx, "key1", 2*time.Second)
	assert.NoError(t, err)
	assert.NotEmpty(t, lockID2)

	err = l.Unlock("key1", lockID2)
	assert.NoError(t, err)

}

// /home/bearbite/.cache/JetBrains/GoLand2024.1/tmp/GoLand/___BenchmarkLockUnlock_in_github_com_trendizz_hydra_hydra_lock.test -test.v -test.paniconexit0 -test.bench ^\QBenchmarkLockUnlock\E$ -test.run ^$
// goos: linux
// goarch: amd64
// pkg: github.com/hydraide/hydraidecore/hydra/lock
// cpu: AMD Ryzen 9 5950X 16-Core Processor
// BenchmarkLockUnlock
// BenchmarkLockUnlock-32    	  533826	      3380 ns/op
// PASS
func BenchmarkLockUnlock(b *testing.B) {
	l := New()

	key := "benchmark_key"

	for i := 0; i < b.N; i++ {

		func() {
			// zárjuk le gyorsan a kontextust minél előbb és ezzel szabadítsuk fel a lockot
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			lockID, err := l.Lock(ctx, key, 2*time.Second)
			if err != nil {
				b.Fatalf("failed to acquire lock: %v", err)
			}
			err = l.Unlock(key, lockID)
			if err != nil {
				b.Fatalf("failed to release lock: %v", err)
			}
		}()

	}

}

// /home/bearbite/.cache/JetBrains/GoLand2024.1/tmp/GoLand/___BenchmarkParallelLockUnlock_in_github_com_trendizz_hydra_hydra_lock.test -test.v -test.paniconexit0 -test.bench ^\QBenchmarkParallelLockUnlock\E$ -test.run ^$
// goos: linux
// goarch: amd64
// pkg: github.com/hydraide/hydraidecore/hydra/lock
// cpu: AMD Ryzen 9 5950X 16-Core Processor
// BenchmarkParallelLockUnlock
// BenchmarkParallelLockUnlock-32    	    8298	    281452 ns/op
// PASS
func BenchmarkParallelLockUnlock(b *testing.B) {
	l := New()
	keys := []string{"key1", "key2", "key3", "key4", "key5"}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for _, key := range keys {

				func() {

					ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
					defer cancel()

					lockID, err := l.Lock(ctx, key, 2*time.Second)
					if err != nil && err.Error() != "lock timeout" {
						b.Fatalf("failed to acquire lock: %v", err)
					}

					if err == nil {
						err = l.Unlock(key, lockID)
						if err != nil {
							b.Fatalf("failed to release lock: %v", err)
						}
					}

				}()

			}
		}
	})
}
