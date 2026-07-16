// Package option provides a generic optional value type.
//
// An Option contains either a value, represented by [Some], or no value,
// represented by [None]. Its zero value is None. An Option can also be created
// from a value and boolean pair using [From], which preserves both parts of the
// pair.
//
// Operations such as [Option.Map], [Option.FlatMap], and [Option.Filter] only
// call their function argument when the Option contains a value. A None
// propagates through these operations without invoking the function.
package option
