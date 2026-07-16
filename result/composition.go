package result

// Fold returns the result of onOK when r is successful and the result of onErr
// otherwise. It calls exactly one of onOK and onErr.
func (r Result[T]) Fold[U any](onOK func(T) U, onErr func(error) U) U {
	if r.err != nil {
		return onErr(r.err)
	}

	return onOK(r.value)
}

// MapErr returns r unchanged when r is successful. When r is failed, MapErr
// calls transform with the contained error and returns a failed Result
// containing the transformed error. Transform must return a non-nil error.
func (r Result[T]) MapErr(transform func(error) error) Result[T] {
	if r.err == nil {
		return r
	}

	return Err[T](transform(r.err))
}

// ZipWith combines r and other by calling combine with both contained values.
// It returns the error from r when r is failed. Otherwise, it returns the error
// from other when other is failed. Combine is called only when both Results are
// successful.
//
// ZipWith is preferred over a Zip method until the module defines a shared
// Pair or Tuple type.
func (r Result[T]) ZipWith[U, V any](other Result[U], combine func(T, U) V) Result[V] {
	if r.err != nil {
		return Err[V](r.err)
	}

	if other.err != nil {
		return Err[V](other.err)
	}

	return Ok(combine(r.value, other.value))
}

// Or returns r when r is successful and other otherwise.
//
// Unlike OrGet, other is evaluated before Or is called.
func (r Result[T]) Or(other Result[T]) Result[T] {
	if r.err != nil {
		return other
	}

	return r
}

// OrGet returns r when r is successful. If r is failed, it calls fallback with
// the contained error and returns the resulting Result. Fallback is not called
// when r is successful.
func (r Result[T]) OrGet(fallback func(error) Result[T]) Result[T] {
	if r.err != nil {
		return fallback(r.err)
	}

	return r
}

// Inspect calls action with the contained value when r is successful, then
// returns r unchanged. It does not call action when r is failed.
func (r Result[T]) Inspect(action func(T)) Result[T] {
	if r.err == nil {
		action(r.value)
	}

	return r
}

// InspectErr calls action with the contained error when r is failed, then
// returns r unchanged. It does not call action when r is successful.
func (r Result[T]) InspectErr(action func(error)) Result[T] {
	if r.err != nil {
		action(r.err)
	}

	return r
}

// Flatten converts a nested Result into a single Result. It returns the outer
// error when the outer Result is failed and otherwise returns the inner Result.
func Flatten[T any](nested Result[Result[T]]) Result[T] {
	if nested.err != nil {
		return Err[T](nested.err)
	}

	return nested.value
}
