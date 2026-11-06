package infcalc

import "math"

/*
I(x). Количество информации (в битах)
(self-information / surprisal / shanon information)
для события появления сообщения idx массива вероятностей probs
*/
func SelfInformation(probs []float64, idx int) float64 {
	if probs[idx] <= 0 {
		return 0
	}
	return -math.Log2(probs[idx])
}

/* H(X). Средняя энтропия источника */
func Entropy(probs []float64) float64 {
	var sum float64
	for i := range probs {
		if probs[i] <= 0 {
			continue
		}
		sum += probs[i] * SelfInformation(probs, i)
	}
	return sum
}

/*
H_max(X) Максимальная энтропия для n сообщений
*/
func MaxEntropy(n int) float64 {
	if n <= 1 {
		return 0
	}
	return math.Log2(float64(n))
}
