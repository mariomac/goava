package errs

import "errors"

// As wraps errors.As without need to first declaring a variable.
// As finds the first error in err's tree that matches target, and if one is found, sets target to that error value and returns true. Otherwise, it returns false.
// The tree consists of err itself, followed by the errors obtained by repeatedly calling its Unwrap() error or Unwrap() []error method. When err wraps multiple errors, As examines err followed by a depth-first traversal of its children.
// An error matches target if the error's concrete value is assignable to the value pointed to by target, or if the error has a method As(any) bool such that As(target) returns true. In the latter case, the As method is responsible for setting target.
// An error type might provide an As method so it can be treated as if it were a different error type.
// As panics if target is not a non-nil pointer to either a type that implements error, or to any interface type.
func As[T error](err error) (T, bool) {
	var v T
	if errors.As(err, &v) {
		return v, true
	}
	return v, false
}
