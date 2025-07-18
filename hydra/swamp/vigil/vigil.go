// Package vigil is a simple synchronization tool that helps monitor whether there is any ongoing operation on
// the bridge in the swamp. Elysium needs this when the graceful stop command has to wait for the completion
// of the last database operation before shutting down Hydra and all other services.
// At the beginning of each database operation, we signal the presence of an active operation by calling the BeginVigil() method.
// At the completion of each database operation, we signal that there are no more active operations by calling the CeaseVigil() method.
package vigil

import (
	"sync"
	"sync/atomic"
)

// Vigil is an interface for managing the state of ongoing operations within the Hydra database.
// It ensures that the database remains operational as long as there are active operations to monitor,
// such as reads and writes. By providing an operation count mechanism, it acts as a safeguard against
// accidentally terminating the database, thus preventing data corruption or loss.
type Vigil interface {

	// BeginVigil indicates that a new operation has started.
	// !!! IMPORTANT !!!: Always call this method before executing any database-related operations.
	// When this method is called, it means that the Hydra database is actively being used
	// and should not be stopped. Each BeginVigil call should be paired with a corresponding CeaseVigil call.
	// Failure to call CeaseVigil will prevent the Hydra database from ever being properly shut down.
	BeginVigil()

	// CeaseVigil indicates that an operation has been completed.
	// When this method is called, it decrements the count of active operations.
	// If there are no more active operations, the Hydra database can be safely shut down.
	CeaseVigil()

	// HasActiveVigils returns a boolean value that indicates whether there are ongoing operations that
	// are currently being monitored. Returns true if active operations exist, false otherwise.
	//
	// Important Note: The function should not be invoked by the Hydra Head because it's the responsibility of the
	// Hydra to determine whether it can be shut down or not using this function!
	HasActiveVigils() bool

	// WaitForActiveVigilsClosed blocks the calling goroutine until all active operations are complete.
	// This ensures that you do not terminate the Hydra database while it's being used, preventing
	// potential data corruption or loss.
	//
	// Important Note: The function should not be invoked by the Hydra Head because it's the responsibility of the
	// Hydra to determine whether it can be shut down or not using this function!
	WaitForActiveVigilsClosed()
}

type vigil struct {
	mu     sync.RWMutex
	cond   *sync.Cond
	vigils int64
}

func New() Vigil {
	v := &vigil{}
	v.cond = sync.NewCond(&v.mu)
	return v
}

func (v *vigil) BeginVigil() {
	atomic.AddInt64(&v.vigils, 1)
}

func (v *vigil) CeaseVigil() {
	atomic.AddInt64(&v.vigils, -1)
	v.cond.Broadcast()
}

func (v *vigil) HasActiveVigils() bool {
	return atomic.LoadInt64(&v.vigils) > 0
}

func (v *vigil) WaitForActiveVigilsClosed() {
	v.cond.L.Lock()
	defer v.cond.L.Unlock()
	for v.HasActiveVigils() {
		v.cond.Wait()
	}
}
