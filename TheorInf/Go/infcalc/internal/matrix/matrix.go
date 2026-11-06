package matrix

import "errors"

func ConcatH[T any](a, b [][]T) ([][]T, error) {
	if len(a) != len(b) {
		return nil, errors.New("mismatched matrix sizes")
	}
	c := make([][]T, len(a))
	for i := range c {
		c[i] = append(a[i], b[i]...)
	}
	return c, nil
}

func E(size int) [][]float64 {
	mx := make([][]float64, size)
	for i := range mx {
		mx[i] = make([]float64, size)
		mx[i][i] = 1
	}
	return mx
}

func Fill[T any](s []T, value T) {
	for i := range s {
		s[i] = value
	}
}

func Contains[T comparable](S [][]T, E []T) bool {
	for i := range S {
		eq := true
		for j := range S[i] {
			if S[i][j] != E[j] {
				eq = false
				break
			}
		}
		if eq {
			return true
		}
	}
	return false
}
