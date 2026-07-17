package seq

import "github.com/everettxwzhu/go-toolkit/option"

// FilterOption lazily applies transform to each value of s and yields the
// contained values from the Options for which transform returns Some.
func (s Seq[T]) FilterOption[U any](transform func(T) option.Option[U]) Seq[U] {
	return FromSeq(func(yield func(U) bool) {
		for value := range s.All() {
			transformed, ok := transform(value).Get()
			if !ok {
				continue
			}

			if !yield(transformed) {
				return
			}
		}
	})
}

// CollectOption applies transform to each value of s and collects the contained
// values. It stops and returns None when transform returns None.
func (s Seq[T]) CollectOption[U any](transform func(T) option.Option[U]) option.Option[[]U] {
	result := make([]U, 0)

	for v := range s.All() {
		transformed, ok := transform(v).Get()
		if !ok {
			return option.None[[]U]()
		}
		result = append(result, transformed)
	}

	return option.Some(result)
}

// FirstOption returns Some containing the first value of s. It returns None
// when s is empty.
func (s Seq[T]) FirstOption() option.Option[T] {
	return option.From(s.First())
}

// LastOption returns Some containing the last value of s. It returns None when
// s is empty.
func (s Seq[T]) LastOption() option.Option[T] {
	return option.From(s.Last())
}

// FindOption returns Some containing the first value for which predicate
// returns true. It returns None when no value matches.
func (s Seq[T]) FindOption(predicate func(T) bool) option.Option[T] {
	return option.From(s.Find(predicate))
}
