package error

import "errors"

var (
	// ErrBinaryValueOutOfByteRange indicates binary value out of byte range.
	ErrBinaryValueOutOfByteRange = errors.New("binary value out of byte range")

	// ErrEmptyProbability indicates empty probability.
	ErrEmptyProbability = errors.New("empty probability")

	// ErrInvalidProbabilityTotal indicates invalid weights or a total other than 100 percent.
	ErrInvalidProbabilityTotal = errors.New("invalid probability total")

	// ErrProbabilityResultNotFound indicates probability result not found.
	ErrProbabilityResultNotFound = errors.New("probability result not found")
)
