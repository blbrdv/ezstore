package utils

import "iter"

func ToSlice[T any](iter iter.Seq[T]) []T {
	var result []T

	for value := range iter {
		result = append(result, value)
	}

	return result
}
