package error

import "errors"

var (
	// ErrBadEmailClient indicates bad email client.
	ErrBadEmailClient = errors.New("bad email client")

	// ErrEmptyToAddresses indicates empty to addresses.
	ErrEmptyToAddresses = errors.New("empty to addresses")

	// ErrBadEmailContent indicates bad email content.
	ErrBadEmailContent = errors.New("bad email content")

	// ErrCreateEmailClient indicates failed to create email client.
	ErrCreateEmailClient = errors.New("failed to create email client")
)
