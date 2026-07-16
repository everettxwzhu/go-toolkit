package result

// Result contains a value of type T and an error.
// It is successful when the error is nil. Its zero value is successful and
// contains the zero value of T.
type Result[T any] struct {
	value T
	err   error
}

// Ok returns a successful Result containing value.
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value}
}

// Err returns a failed Result containing err.
// It panics if err is nil.
func Err[T any](err error) Result[T] {
	if err == nil {
		panic("result: Err called with nil")
	}

	return Result[T]{err: err}
}

// From returns a Result containing value and err.
// The Result is successful when err is nil and failed otherwise.
func From[T any](value T, err error) Result[T] {
	return Result[T]{
		value: value,
		err:   err,
	}
}

// Get returns the value and error contained in r.
func (r Result[T]) Get() (T, error) {
	return r.value, r.err
}

// IsOK reports whether r is successful.
func (r Result[T]) IsOK() bool {
	return r.err == nil
}

// IsErr reports whether r is failed.
func (r Result[T]) IsErr() bool {
	return r.err != nil
}

// OrElse returns the contained value when r is successful and fallback
// otherwise.
func (r Result[T]) OrElse(fallback T) T {
	if r.err != nil {
		return fallback
	}

	return r.value
}

// OrElseGet returns the contained value when r is successful. If r is failed,
// it calls fallback with the contained error and returns the resulting value.
func (r Result[T]) OrElseGet(fallback func(error) T) T {
	if r.err != nil {
		return fallback(r.err)
	}

	return r.value
}

// Map returns a Result produced by applying transform to the contained value.
// If r is failed, Map propagates its error without calling transform.
func (r Result[T]) Map[U any](transform func(T) U) Result[U] {
	if r.err != nil {
		return Err[U](r.err)
	}

	return Ok(transform(r.value))
}

// FlatMap returns the Result produced by applying transform to the contained
// value. If r is failed, FlatMap propagates its error without calling
// transform.
func (r Result[T]) FlatMap[U any](transform func(T) Result[U]) Result[U] {
	if r.err != nil {
		return Err[U](r.err)
	}

	return transform(r.value)
}
