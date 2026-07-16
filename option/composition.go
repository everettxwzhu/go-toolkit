package option

import "github.com/everettxwzhu/go-toolkit/result"

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
//
// ZipWith is preferred over a Zip method until the module defines a shared
// Pair or Tuple type.
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
