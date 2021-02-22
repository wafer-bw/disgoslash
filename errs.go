package disgoslash

import "errors"

// ErrUnauthorized - the request signature is invalid or Discord API responded with 401
var ErrUnauthorized = errors.New("unauthorized")

// ErrInvalidInteractionType - the request interaction type is invalid
var ErrInvalidInteractionType = errors.New("invalid interaction type")

// ErrNotImplemented - whatever was requested is not implemented yet
var ErrNotImplemented = errors.New("not implemented")

// ErrAlreadyExists - returned when attempting to create a command which already exists
var ErrAlreadyExists = errors.New("already exists")

// ErrTooManyRequests - returned when the Disord API responds with a 429
var ErrTooManyRequests = errors.New("too many requests")

// ErrForbidden - returned when the Disord API responds with a 403
var ErrForbidden = errors.New("forbidden - missing access")

// ErrMaxRetries - returned when the maximum number of retries is reached in a retry loop
var ErrMaxRetries = errors.New("max retries reached")
