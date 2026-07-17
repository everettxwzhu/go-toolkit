package option

import (
	"github.com/everettxwzhu/go-toolkit/result"
	"github.com/everettxwzhu/go-toolkit/tuple"
)

// Fold returns the result of onSome when o contains a value and the result of
// onNone otherwise. It calls exactly one of onSome and onNone.
func (o Option[T]) Fold[U any](onSome func(T) U, onNone func() U) U {
	if !o.ok {
		return onNone()
	}

	return onSome(o.value)
}

// ZipWith combines o and other by calling combine with both contained values.
// It returns None without calling combine when either Option is None.
func (o Option[T]) ZipWith[U, V any](other Option[U], combine func(T, U) V) Option[V] {
	if !o.ok || !other.ok {
		return None[V]()
	}

	return Some(combine(o.value, other.value))
}

// Or returns o when o is Some and other otherwise.
//
// Unlike OrGet, other is evaluated before Or is called.
func (o Option[T]) Or(other Option[T]) Option[T] {
	if !o.ok {
		return other
	}

	return o
}

// OrGet returns o when o is Some. If o is None, it calls fallback and returns
// the resulting Option. Fallback is not called when o is Some.
func (o Option[T]) OrGet(fallback func() Option[T]) Option[T] {
	if !o.ok {
		return fallback()
	}

	return o
}

// Inspect calls action with the contained value when o is Some, then returns o
// unchanged. It does not call action when o is None.
func (o Option[T]) Inspect(action func(T)) Option[T] {
	if o.ok {
		action(o.value)
	}

	return o
}

// Flatten converts a nested Option into a single Option. It returns None when
// either the outer or inner Option is None.
func Flatten[T any](nested Option[Option[T]]) Option[T] {
	if !nested.ok {
		return None[T]()
	}

	if !nested.value.ok {
		return None[T]()
	}

	return Some(nested.value.value)
}

// ToResult returns a successful result containing the value when o is Some.
// When o is None, it returns a failed result containing err. Err must be
// non-nil when o is None.
//
// Keep this conversion in package option so only option imports result; package
// result must not import option, which would create an import cycle.
func (o Option[T]) ToResult(err error) result.Result[T] {
	if !o.ok {
		return result.Err[T](err)
	}

	return result.Ok(o.value)
}

// ToResultGet returns a successful result containing the value when o is Some.
// When o is None, it calls err and returns a failed result containing the
// returned error. Err is not called when o is Some and must return a non-nil
// error when called.
func (o Option[T]) ToResultGet(err func() error) result.Result[T] {
	if !o.ok {
		return result.Err[T](err())
	}

	return result.Ok(o.value)
}

// FromResult returns Some containing the value when r is successful and None
// when r is failed. The error in a failed Result is discarded.
func FromResult[T any](r result.Result[T]) Option[T] {
	if r.IsErr() {
		return None[T]()
	}

	v, _ := r.Get()
	return Some(v)
}

// Zip pairs the contained values from o and other. It returns None when either
// Option is None.
func Zip[T, U any](left Option[T], right Option[U]) Option[tuple.Pair[T, U]] {
	if left.IsNone() || right.IsNone() {
		return None[tuple.Pair[T, U]]()
	}

	return Some(tuple.New(left.value, right.value))
}

// Exists reports whether o is Some and predicate returns true for its contained
// value. Predicate is not called when o is None.
func (o Option[T]) Exists(predicate func(T) bool) bool {
	if o.IsNone() {
		return false
	}

	return predicate(o.value)
}

// ContainsBy reports whether o is Some and its contained value is equal to
// target according to equal.
func (o Option[T]) ContainsBy(target T, equal func(T, T) bool) bool {
	if o.IsNone() {
		return false
	}

	return equal(o.value, target)
}

// ToSlice returns a one-element slice containing the value when o is Some and
// an empty slice when o is None.
func (o Option[T]) ToSlice() []T {
	if o.IsNone() {
		return []T{}
	}

	return []T{o.value}
}

// Transpose converts Some(Ok(value)) to Ok(Some(value)), Some(Err(err)) to
// Err(err), and None to Ok(None).
func Transpose[T any](nested Option[result.Result[T]]) result.Result[Option[T]] {
	if nested.IsNone() {
		return result.Ok(None[T]())
	}

	v, err := nested.value.Get()
	if nested.value.IsErr() {
		return result.Err[Option[T]](err)
	}

	return result.Ok(Some(v))
}

// Sequence converts a slice of Options into an Option containing all values.
// It returns None at the first None and otherwise preserves value order.
func Sequence[T any](values []Option[T]) Option[[]T] {
	list := make([]T, 0, len(values))

	for _, v := range values {
		if v.IsNone() {
			return None[[]T]()
		}
		list = append(list, v.value)
	}

	return Some(list)
}

// Traverse applies transform to each value in order and collects the contained
// values. It stops and returns None when transform returns None.
func Traverse[T, U any](values []T, transform func(T) Option[U]) Option[[]U] {
	list := make([]U, 0, len(values))

	for _, v := range values {
		o := transform(v)
		if o.IsNone() {
			return None[[]U]()
		}
		list = append(list, o.value)
	}

	return Some(list)
}

// Contains reports whether o is Some and contains target.
//
// Contains is a top-level function because a method cannot add a comparable
// constraint to the existing type parameter T of Option[T].
func Contains[T comparable](o Option[T], target T) bool {
	if o.IsNone() {
		return false
	}

	return o.value == target
}
