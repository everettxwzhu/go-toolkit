package result

import (
	"fmt"

	"github.com/everettxwzhu/go-toolkit/tuple"
)

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

// Zip pairs the contained values from left and right. It returns the error from
// left when left is failed, otherwise the error from right when right is failed.
func Zip[T, U any](left Result[T], right Result[U]) Result[tuple.Pair[T, U]] {
	if left.IsErr() {
		return Err[tuple.Pair[T, U]](left.err)
	}
	if right.IsErr() {
		return Err[tuple.Pair[T, U]](right.err)
	}

	return Ok(tuple.New(left.value, right.value))
}

// Ensure returns r unchanged when r is failed or predicate returns true for its
// contained value. When predicate returns false, Ensure calls err with the
// value and returns a failed Result containing the returned error.
//
// Predicate and err are not called when r is already failed. Err must return a
// non-nil error when called.
func (r Result[T]) Ensure(predicate func(T) bool, err func(T) error) Result[T] {
	if r.IsErr() {
		return r
	}

	if predicate(r.value) {
		return r
	}

	return Err[T](err(r.value))
}

// ContainsBy reports whether r is successful and its contained value is equal
// to target according to equal.
func (r Result[T]) ContainsBy(target T, equal func(T, T) bool) bool {
	if r.IsErr() {
		return false
	}

	return equal(r.value, target)
}

// Sequence converts a slice of Results into a Result containing all successful
// values. It stops at the first failed Result and otherwise preserves value
// order.
func Sequence[T any](values []Result[T]) Result[[]T] {
	list := make([]T, 0, len(values))

	for _, v := range values {
		if v.IsErr() {
			return Err[[]T](v.err)
		}
		list = append(list, v.value)
	}

	return Ok(list)
}

// Traverse applies transform to each value in order and collects the successful
// values. It stops and returns the first failed Result.
func Traverse[T, U any](values []T, transform func(T) Result[U]) Result[[]U] {
	list := make([]U, 0, len(values))

	for _, v := range values {
		o := transform(v)
		if o.IsErr() {
			return Err[[]U](o.err)
		}
		list = append(list, o.value)
	}

	return Ok(list)
}

// WrapErr returns r unchanged when r is successful. When r is failed, WrapErr
// returns a failed Result whose error adds message while preserving the
// original error for errors.Is and errors.As.
func (r Result[T]) WrapErr(message string) Result[T] {
	if r.IsErr() {
		return Err[T](fmt.Errorf("%s: %w", message, r.err))
	}

	return r
}
