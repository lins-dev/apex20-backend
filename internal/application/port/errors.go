package port

import "errors"

// ErrNotFound is returned by repository methods when the requested record does not exist.
var ErrNotFound = errors.New("not found")

// ErrEmailAlreadyExists is returned when trying to register an e-mail that is already in use.
var ErrEmailAlreadyExists = errors.New("email already exists")

// ErrInvalidCredentials is returned when the provided password does not match the stored hash.
var ErrInvalidCredentials = errors.New("invalid credentials")
