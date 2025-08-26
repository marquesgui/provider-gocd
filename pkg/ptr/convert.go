// Package ptr provides generic helpers for working with pointers.
package ptr

func ToPtr[T any](v T) *T {
	return &v
}
