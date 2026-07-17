package seq

import (
	"github.com/everettxwzhu/go-toolkit/option"
	"github.com/everettxwzhu/go-toolkit/result"
)

// TryCollect applies transform to each value of s and collects the transformed
// values. It stops at the first error and returns that error without retaining
// the values collected before the error.
func (s Seq[T]) TryCollect[U any](transform func(T) (U, error)) result.Result[[]U] {
	values := make([]U, 0)

	for t := range s.All() {
		if u, err := transform(t); err != nil {
			return result.Err[[]U](err)
		} else {
			values = append(values, u)
		}
	}

	return result.Ok(values)
}

// TryReduce combines initial and the values of s from left to right. It stops
// at the first error returned by reducer and returns that error.
func (s Seq[T]) TryReduce[U any](initial U, reducer func(U, T) (U, error)) result.Result[U] {
	value := initial

	for t := range s.All() {
		if v, err := reducer(value, t); err != nil {
			return result.Err[U](err)
		} else {
			value = v
		}
	}

	return result.Ok(value)
}

// TryForEach calls action once for each value of s in iteration order. It stops
// at the first error and returns it.
func (s Seq[T]) TryForEach(action func(T) error) error {
	for t := range s.All() {
		if err := action(t); err != nil {
			return err
		}
	}

	return nil
}

// TryAny reports whether predicate returns true for at least one value of s. It
// stops at the first match or the first error.
func (s Seq[T]) TryAny(predicate func(T) (bool, error)) result.Result[bool] {
	for t := range s.All() {
		if match, err := predicate(t); err != nil {
			return result.Err[bool](err)
		} else if match {
			return result.Ok(true)
		}
	}

	return result.Ok(false)
}

// TryEvery reports whether predicate returns true for every value of s. It
// stops at the first false result or the first error.
func (s Seq[T]) TryEvery(predicate func(T) (bool, error)) result.Result[bool] {
	for t := range s.All() {
		if match, err := predicate(t); err != nil {
			return result.Err[bool](err)
		} else if !match {
			return result.Ok(false)
		}
	}

	return result.Ok(true)
}

// TryFind returns the first value for which predicate returns true. It returns
// None when no value matches and stops at the first error.
func (s Seq[T]) TryFind(predicate func(T) (bool, error)) result.Result[option.Option[T]] {
	for t := range s.All() {
		if match, err := predicate(t); err != nil {
			return result.Err[option.Option[T]](err)
		} else if match {
			return result.Ok(option.Some(t))
		}
	}

	return result.Ok(option.None[T]())
}

// MapResult lazily applies transform to each value of s and wraps each value or
// error in an independent Result. Unlike TryCollect, MapResult does not stop
// after a failed transformation.
func (s Seq[T]) MapResult[U any](transform func(T) (U, error)) Seq[result.Result[U]] {
	return FromSeq(func(yield func(result.Result[U]) bool) {
		for t := range s.All() {
			transformed, err := transform(t)
			r := result.From(transformed, err)

			if !yield(r) {
				return
			}
		}
	})
}
