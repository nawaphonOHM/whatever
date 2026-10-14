// Package config provides configuration parsing, validation, and defaults
// for object storage connections.
package config

import (
	"context"
	"fmt"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/central"
)

// SetExitFunc overrides the process exit hook used by the central logger worker
// and returns the previous hook.
func SetExitFunc(fn func(int)) func(int) {
	// False positive: borrowing the central worker singleton to configure exit hook is leak-free;
	// proof: TestDefaultWorker_SingletonResourceProof in worker_singleton_proof_test.go.
	prev := central.DefaultWorker().ExitFunc()
	central.DefaultWorker().SetExitFunc(fn)
	return prev
}

// handleMissingKeys inspects missing keys and initiates peaceful termination via Central Log if any are missing.
func handleMissingKeys(missing []string) error {
	if len(missing) == 0 {
		return nil
	}
	msg := fmt.Sprintf("missing required storage configuration keys: %v; exiting peacefully", missing)
	central.ExitWithGraceful(context.Background(), msg)
	return fmt.Errorf("missing required storage configuration keys: %v", missing)
}
