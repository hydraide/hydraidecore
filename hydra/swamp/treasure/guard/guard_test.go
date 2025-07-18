package guard

import (
	"sync"
	"testing"
)

// MyObject is a sample object that uses the Guard interface to ensure
// thread-safe access to its parameters.
type MyObject struct {
	Guard
	firstParam  int
	secondParam int
}

// NewMyObject initializes a new MyObject with a Guard.
func NewMyObject() *MyObject {
	return &MyObject{
		Guard: New(),
	}
}

// IncreaseFirstParam safely increments the firstParam field.
// It uses the Guard to ensure that only one goroutine can execute this method at a time.
func (o *MyObject) IncreaseFirstParam(lockerID ID) {
	_ = o.Guard.CanExecute(lockerID)
	o.firstParam++
}

// IncreaseSecondParam safely increments the secondParam field.
// It uses the Guard to ensure that only one goroutine can execute this method at a time.
func (o *MyObject) IncreaseSecondParam(lockerID ID) {
	_ = o.Guard.CanExecute(lockerID)
	o.secondParam++
}

// TestNew tests the NewMyObject function and the Guard implementation.
// It spawns 5000 goroutines that try to increment the parameters of a shared MyObject instance.
func TestNew(t *testing.T) {

	myObject := NewMyObject()

	wg := sync.WaitGroup{}
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(i int, myObject *MyObject) {
			lockerID := myObject.StartTreasureGuard(true)
			myObject.IncreaseFirstParam(lockerID)
			myObject.IncreaseSecondParam(lockerID)
			myObject.ReleaseTreasureGuard(lockerID)
			wg.Done()
		}(i, myObject)
	}

	wg.Wait()
}

// BenchmarkNew benchmarks the Guard implementation.
// It measures the time it takes to acquire and release a transaction using the Guard.
// goos: windows
// goarch: amd64
// cpu: AMD Ryzen 9 5950X 16-Core Processor
//
// BenchmarkNew
// BenchmarkNew-26         28281244                41.98 ns/op
// PASS
//
// goos: linux
// goarch: amd64
// pkg: github.com/hydraide/hydraidecore/hydra/swamp/treasure/guard
// cpu: AMD Ryzen 9 5950X 16-Core Processor
//
// BenchmarkNew
// BenchmarkNew-32    	30823390	        54.59 ns/op
// goos: linux
// goarch: amd64
// pkg: github.com/hydraide/hydraidecore/hydra/swamp/treasure/guard
// cpu: AMD Ryzen 9 5950X 16-Core Processor
// BenchmarkNew
// BenchmarkNew-32    	24532616	        66.54 ns/op
// PASS
func BenchmarkNew(b *testing.B) {
	myObject := NewMyObject()
	for i := 0; i < b.N; i++ {
		lockerID := myObject.StartTreasureGuard(true)
		if lockerID == 0 {
			continue
		} else {
			_ = myObject.CanExecute(lockerID)
			myObject.ReleaseTreasureGuard(lockerID)
		}
	}
}
