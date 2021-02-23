package disgoslash

import "errors"

// ErrUnauthorized is returned when the request signature is invalid or Discord API responded with 401
var ErrUnauthorized = errors.New("unauthorized")

// ErrInvalidInteractionType is returned when the request interaction type is invalid
var ErrInvalidInteractionType = errors.New("invalid interaction type")

// ErrNotImplemented is returned when whatever was requested hasn't been implemented yet
var ErrNotImplemented = errors.New("not implemented")

// ErrAlreadyExists is returned when attempting to create a command which already exists
var ErrAlreadyExists = errors.New("already exists")

// ErrTooManyRequests is returned when the Disord API responds with a 429
var ErrTooManyRequests = errors.New("too many requests")

// ErrForbidden is returned when the Disord API responds with a 403
var ErrForbidden = errors.New("forbidden - missing access")

// ErrMaxRetries is returned when the maximum number of retries is reached in a retry loop
var ErrMaxRetries = errors.New("max retries reached")
