package option

// Option contains either a value of type T or no value.
// Its zero value is None.
type Option[T any] struct {
	value T
	ok    bool
}

// Some returns an Option containing value.
func Some[T any](value T) Option[T] {
	return Option[T]{
		value: value,
		ok:    true,
	}
}

// None returns an Option containing no value.
func None[T any]() Option[T] {
	return Option[T]{}
}

// From returns an Option containing value with the presence state ok.
// Both value and ok are preserved as supplied.
func From[T any](value T, ok bool) Option[T] {
	return Option[T]{
		value: value,
		ok:    ok,
	}
}

// Get returns the value and presence state contained in o.
func (o Option[T]) Get() (T, bool) {
	return o.value, o.ok
}

// IsSome reports whether o contains a value.
func (o Option[T]) IsSome() bool {
	return o.ok
}

// IsNone reports whether o contains no value.
func (o Option[T]) IsNone() bool {
	return !o.ok
}

// OrElse returns the contained value when o is Some and fallback otherwise.
func (o Option[T]) OrElse(fallback T) T {
	if !o.ok {
		return fallback
	}

	return o.value
}

// OrElseGet returns the contained value when o is Some. If o is None, it calls
// fallback and returns the resulting value.
func (o Option[T]) OrElseGet(fallback func() T) T {
	if !o.ok {
		return fallback()
	}

	return o.value
}

// Map returns an Option produced by applying transform to the contained value.
// If o is None, Map returns None without calling transform.
func (o Option[T]) Map[U any](transform func(T) U) Option[U] {
	if !o.ok {
		return None[U]()
	}

	return Some(transform(o.value))
}

// FlatMap returns the Option produced by applying transform to the contained
// value. If o is None, FlatMap returns None without calling transform.
func (o Option[T]) FlatMap[U any](transform func(T) Option[U]) Option[U] {
	if !o.ok {
		return None[U]()
	}

	return transform(o.value)
}

// Filter returns o when it is Some and predicate returns true. It returns None
// when o is None or predicate returns false. Predicate is not called for None.
func (o Option[T]) Filter(predicate func(T) bool) Option[T] {
	if !o.ok {
		return None[T]()
	}

	if !predicate(o.value) {
		return None[T]()
	}

	return o
}
