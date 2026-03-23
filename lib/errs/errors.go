package errs

import "errors"

// ErrForbidden is returned when the caller lacks permission for the requested operation.
var ErrForbidden = errors.New("forbidden")
