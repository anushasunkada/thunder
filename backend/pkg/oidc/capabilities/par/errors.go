package par

import "errors"

var (
	// ErrRequestURINotFound is returned when the request_uri is unknown or has expired.
	ErrRequestURINotFound = errors.New("par: request_uri not found or expired")

	// ErrClientIDMismatch is returned when the client_id in the authorize request
	// does not match the one stored in the pushed authorization request.
	ErrClientIDMismatch = errors.New("par: client_id mismatch")
)
