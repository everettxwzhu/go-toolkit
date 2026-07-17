// Package tuple provides small generic product types used when an
// operation naturally returns multiple differently typed values.
package tuple

// Pair is an ordered pair containing values of types A and B.
type Pair[A, B any] struct {
	First  A
	Second B
}

// New returns a Pair containing first and second.
func New[A, B any](first A, second B) Pair[A, B] {
	return Pair[A, B]{
		First:  first,
		Second: second,
	}
}

// Unpack returns the two values contained in p.
func (p Pair[A, B]) Unpack() (A, B) {
	return p.First, p.Second
}

// Swap returns a Pair with the values of p exchanged.
func (p Pair[A, B]) Swap() Pair[B, A] {
	return Pair[B, A]{
		First:  p.Second,
		Second: p.First,
	}
}

// MapFirst returns a Pair produced by applying transform to the first value of
// p while preserving the second value.
func (p Pair[A, B]) MapFirst[C any](transform func(A) C) Pair[C, B] {
	first := transform(p.First)
	return Pair[C, B]{
		First:  first,
		Second: p.Second,
	}
}

// MapSecond returns a Pair produced by applying transform to the second value
// of p while preserving the first value.
func (p Pair[A, B]) MapSecond[C any](transform func(B) C) Pair[A, C] {
	second := transform(p.Second)
	return Pair[A, C]{
		First:  p.First,
		Second: second,
	}
}
