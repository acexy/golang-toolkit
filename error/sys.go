package error

import "errors"

var (
	// ErrShutdownSignalsInitialized indicates shutdown signals already initialized.
	ErrShutdownSignalsInitialized = errors.New("shutdown signals already initialized")

	// ErrInvalidShutdownSignals indicates shutdown signals cannot be empty or contain nil.
	ErrInvalidShutdownSignals = errors.New("shutdown signals cannot be empty or contain nil")
)
