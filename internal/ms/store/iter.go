package store

import "iter"

// ToSlice converts iterator over sequence to slice.
func ToSlice[T any](iter iter.Seq[T]) []T {
	var result []T

	for value := range iter {
		result = append(result, value)
	}

	return result
}
