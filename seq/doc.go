// Package seq provides a generic, iterator-based sequence type and common
// operations for constructing, transforming, selecting, combining, and
// consuming sequences.
//
// A Seq is lazy: intermediate operations such as [Seq.Map], [Seq.Filter], and
// [Seq.Concat] build a new sequence without traversing their input. Values are
// produced only when the sequence is iterated with [Seq.All] or consumed by a
// terminal operation such as [Seq.Collect], [Seq.Reduce], or [Seq.Count].
// Iteration stops early when the iterator's yield function returns false, and
// short-circuiting terminal operations such as [Seq.First], [Seq.Find], and
// [Seq.Any] request only the values they need.
//
// Sequences can be created from slices, iterator functions, variadic values,
// and integer ranges using [FromSlice], [FromSeq], [Of], and [Range]. The zero
// value of Seq is valid and represents an empty sequence.
//
// Whether a sequence can be iterated more than once depends on its source.
// Operations in this package do not cache produced values. Use [Seq.Collect]
// when a materialized snapshot is required.
package seq
