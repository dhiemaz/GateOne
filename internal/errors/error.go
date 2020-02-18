package errors

type ErrorCode string

// Causer functions list
type Causer interface {
	Cause() error
}

// Error data structure
type Error struct {
	code    ErrorCode
	origin  error
	message string
}

// Error is ...
func (e Error) Error() string {
	return e.message
}

// Cause is ...
func (e Error) Cause() error {
	return e.origin
}

// Code is ...
func (e Error) Code() ErrorCode {
	return e.code
}
