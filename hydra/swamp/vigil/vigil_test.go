package vigil

import (
	"sync/atomic"
	"testing"
)

func TestNew(t *testing.T) {

	vigilObj := New()
	vigilObj.BeginVigil()

	if !vigilObj.HasActiveVigils() {
		t.Errorf("Expected true, got false")
	}

	vigilObj.CeaseVigil()

	if vigilObj.HasActiveVigils() {
		t.Errorf("Expected false, got true")
	}

}

func TestWaitingForUnlock(t *testing.T) {

	vigilObj := New()

	allTransactions := int32(10000)

	counter := int32(0)
	for i := int32(0); i < allTransactions; i++ {
		vigilObj.BeginVigil()
		go func() {
			atomic.AddInt32(&counter, 1)
			vigilObj.CeaseVigil()
		}()
	}

	vigilObj.WaitForActiveVigilsClosed()

	if counter != allTransactions {
		t.Errorf("Expected 100, got %d", counter)
	}

}

// goos: windows
// goarch: amd64
// pkg: github.com/trendizz/neendb/neen/transaction
// cpu: AMD Ryzen 9 5950X 16-Core Processor
// BenchmarkNew
// BenchmarkNew-32         235088299                5.104 ns/op
// PASS
func BenchmarkNew(b *testing.B) {

	vigilObj := New()

	for i := 0; i < b.N; i++ {
		vigilObj.BeginVigil()
		vigilObj.CeaseVigil()
	}

}
