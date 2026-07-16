// Package result provides a generic value-or-error type.
//
// A Result contains a value and an error. It is successful when its error is
// nil and failed otherwise. Its zero value is a successful Result containing
// the zero value of its value type. Results can be created explicitly with
// [Ok] and [Err], or from a conventional value and error pair using [From].
//
// Operations such as [Result.Map] and [Result.FlatMap] only call their function
// argument for a successful Result. A failed Result propagates its error
// through these operations without invoking the function.
package result
