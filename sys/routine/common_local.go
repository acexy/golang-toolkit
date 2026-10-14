package routine

import (
	"encoding/hex"

	"github.com/google/uuid"
)

// NewTraceIdThreadLocal creates inheritable trace storage; a nil supplier generates a 32-character UUID.
// Call Set or Get in the parent before Run or RunAll to share one identifier.
// Reuse the instance and Delete the current value when its request finishes.
func NewTraceIdThreadLocal(supplier Supplier[string]) *InheritableThreadLocal[string] {
	if supplier == nil {
		supplier = func() string {
			id := uuid.New()
			return hex.EncodeToString(id[:])
		}
	}
	return NewInheritableThreadLocal(supplier)
}
