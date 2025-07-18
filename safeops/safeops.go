// Package safeops provides a simple mechanism for locking the system during database-related operations to ensure data integrity
// and prevent the system from shutting down while any operation is running.
package safeops

import (
	"sync"
	"sync/atomic"
	"time"
)

type Safeops interface {

	// LockSystem prevents the system from shutting down. This function must be invoked before executing any database-related operations.
	// Its purpose is to prevent the system from shutting down any database-related processes or the system itself while the database operation is running.
	LockSystem()

	// MonitorPanic returns a channel that can be monitored. If a true message arrives on the channel, it indicates that a panic has occurred,
	// and the safeops emergency shutdown process should begin. During this process, no new requests are accepted; only ongoing requests are completed.
	// Real-world scenario: Useful for handling unexpected crashes or panics in the system, ensuring data integrity before shutting down.
	MonitorPanic() chan bool

	// UnlockSystem serves to unlock the system. It signals to the system that the database-related operation has completed, allowing the system to shut down if it wishes.
	UnlockSystem()

	// SystemLocked returns true if at least one process has requested a system transaction and has not yet released it.
	SystemLocked() bool

	// TriggerPanic is called when a panic occurs in the system. It sends a true value to the stopSignal channel, which can be monitored using WatchForApocalypse.
	TriggerPanic()

	// WaitForUnlock waits until the system releases the transaction. The function returns when there are no more active transaction requests in the system.
	// This is a blocking process with automatic release, which should be called before gracefulStop, allowing us to wait for the system to properly handle the locks.
	// Real-world scenario: Useful for ensuring that all critical operations have completed before initiating a safeops shutdown.
	WaitForUnlock()
}

type safeops struct {
	mu         sync.RWMutex
	stopSignal chan bool
	isLocked   int32
}

func New() Safeops {
	g := &safeops{
		stopSignal: make(chan bool),
	}
	return g
}

func (s *safeops) LockSystem() {
	atomic.AddInt32(&s.isLocked, 1)
}

func (s *safeops) UnlockSystem() {
	atomic.AddInt32(&s.isLocked, -1)
}

func (s *safeops) SystemLocked() bool {
	return atomic.LoadInt32(&s.isLocked) > 0
}

func (s *safeops) MonitorPanic() chan bool {
	return s.stopSignal
}

func (s *safeops) TriggerPanic() {
	s.stopSignal <- true
}

func (s *safeops) WaitForUnlock() {
	// waiting until the s.isLocked will be 0
	for {
		if !s.SystemLocked() {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}
