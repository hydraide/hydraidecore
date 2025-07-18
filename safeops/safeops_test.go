package safeops

import (
	"testing"
	"time"
)

func TestSealTheGates(t *testing.T) {

	e := New()

	// Test to ensure the gates are not sealed initially
	if e.SystemLocked() {
		t.Errorf("The Safeops gates are sealed by default, which should not be the case.")
	}

	// Invoke the LockSystem method
	e.LockSystem()

	// Test to ensure the gates are now sealed
	if !e.SystemLocked() {
		t.Errorf("The Safeops gates are not sealed, but they should be.")
	}

}

func TestUnsealTheGates(t *testing.T) {

	e := New()

	// First, seal the gates to set up the test scenario
	e.LockSystem()

	// Test to ensure the gates are indeed sealed
	if !e.SystemLocked() {
		t.Errorf("The Safeops gates are not sealed, but they should be for this test.")
	}

	// Invoke the UnlockSystem method
	e.UnlockSystem()

	// Test to ensure the gates are now unsealed
	if e.SystemLocked() {
		t.Errorf("The Safeops gates are still sealed, but they should be unsealed.")
	}

}

func TestGatesSealed(t *testing.T) {

	e := New()

	// Test to ensure that the gates are initially unsealed
	if e.SystemLocked() {
		t.Errorf("The Safeops gates are initially sealed, but they should be unsealed.")
	}

	// Seal the gates
	e.LockSystem()

	// Test to ensure the gates are now sealed
	if !e.SystemLocked() {
		t.Errorf("The Safeops gates are not sealed, but they should be sealed.")
	}

	// Unseal the gates
	e.UnlockSystem()

	// Test to ensure the gates are unsealed again
	if e.SystemLocked() {
		t.Errorf("The Safeops gates are still sealed, but they should be unsealed.")
	}
}

func TestWatchForApocalypse(t *testing.T) {

	e := New()

	// Create a channel to capture the panic signal
	apocalypseChan := e.MonitorPanic()

	// Create a channel to signal test completion
	done := make(chan bool)

	go func() {
		// Watch for apocalypse signal
		select {
		case <-apocalypseChan:
			// Signal received, test should pass
			done <- true
		case <-time.After(1 * time.Second):
			// Timeout, test should fail
			done <- false
		}
	}()

	// Sound the horns to trigger apocalypse
	e.TriggerPanic()

	// Wait for the result
	if !<-done {
		t.Errorf("Failed to receive apocalypse signal, MonitorPanic is not working as expected.")
	}
}

func TestAwaitUnsealing(t *testing.T) {

	e := New()

	// Initially, the gates should be unsealed
	if e.SystemLocked() {
		t.Errorf("Initially, the gates should not be sealed.")
	}

	// Create a channel to signal test completion
	done := make(chan bool)

	go func() {
		// Listen for unsealing completion
		e.WaitForUnlock()
		done <- true
	}()

	// Seal the gates
	e.LockSystem()
	if !e.SystemLocked() {
		t.Errorf("After sealing, the gates should be sealed.")
	}

	// Unseal the gates
	e.UnlockSystem()
	if e.SystemLocked() {
		t.Errorf("After unsealing, the gates should not be sealed.")
	}

	// Check if the WaitForUnlock has finished
	select {
	case <-done:
		// Test passed
	case <-time.After(1 * time.Second):
		t.Errorf("WaitForUnlock did not complete as expected.")
	}
}
