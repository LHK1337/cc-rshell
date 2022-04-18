package utils

import "golang.org/x/exp/constraints"

func Max[T constraints.Ordered](v1 T, values ...T) T {
	max := v1

	for _, v := range values[1:] {
		if v > max {
			max = v
		}
	}

	return max
}

func Min[T constraints.Ordered](v1 T, values ...T) T {
	min := v1

	for _, v := range values {
		if v < min {
			min = v
		}
	}

	return min
}
