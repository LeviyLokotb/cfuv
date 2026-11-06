package infcalc

import "math/rand/v2"

/*
P(X). Заполняет массив случайными вероятностями появления
сообщений на входе информационного устройства,
в сумме 1.0f. Возвращает полученную сумму вероятностей
*/
func NewProbsIn(n int) []float64 {
	probs := make([]float64, n)
	var sum float64
	for i := range probs {
		r := rand.Float64()
		probs[i] = r
		sum += r
	}
	if sum == 0 {
		for i := range probs {
			probs[i] = 1.0 / float64(n)
		}
		return probs
	}
	for i := range probs {
		probs[i] /= sum
	}
	return probs
}

/*
P(Y/X). Заполняет матрицу n x n случайными условными вероятностями
получения сообщения на выходе при известном сообщении на входе.

# Первый индекс (строка, i) -- индекс переданного сообщения

# Второй индекс (столбец, j) -- индекс полученного сообщения

Сумма по строке всегда 1.0f. Вероятность безошибочной передачи не менее minCorrectProb
*/
func NewConditionalProbs(n int, minCorrectProb float64) [][]float64 {
	condProbs := make([][]float64, n)
	for i := range condProbs {
		condProbs[i] = make([]float64, n)
	}
	if n == 1 {
		condProbs[0][0] = 1.0
		return condProbs
	}
	for i := 0; i < n; i++ {
		var sum float64
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			condProbs[i][j] = rand.Float64()
			sum += condProbs[i][j]
		}
		incorrectProb := rand.Float64() * (1 - minCorrectProb)
		condProbs[i][i] = 1 - incorrectProb
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			if sum <= 0 {
				condProbs[i][j] = incorrectProb / float64(n-1)
				continue
			}
			condProbs[i][j] /= sum
			condProbs[i][j] *= incorrectProb
		}
	}
	return condProbs
}

/*
P(Y). Заполняет массив вероятностями появления
сообщений на выходе информационного устройства.

probsIn -- массив вероятностей появления сообщений на входе P(X)

condProbs -- матрица условных вероятностей P(Y/X)
*/
func NewProbsOut(n int, probsIn []float64, condProbs [][]float64) []float64 {
	probsOut := make([]float64, n)
	for j := 0; j < n; j++ {
		var sum float64
		for i := 0; i < n; i++ {
			sum += probsIn[i] * condProbs[i][j]
		}
		probsOut[j] = sum
	}
	return probsOut
}

/*
P(X,Y). Заполняет матрицу n x n совместными вероятностями
появления конкретных сообщений на входе и на выходе

probsIn -- массив вероятностей появления сообщений на входе P(X)

condProbs -- матрица условных вероятностей P(Y/X)
*/
func NewJointProbs(n int, probsIn []float64, condProbs [][]float64) [][]float64 {
	jointProbs := make([][]float64, n)
	for i := range jointProbs {
		jointProbs[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			jointProbs[i][j] = probsIn[i] * condProbs[i][j]
		}
	}
	return jointProbs
}

/*
τ. Заполняет массив случайными длительностями символа
в диапазоне (0, n].
*/
func NewDurations(n int) []float64 {
	durations := make([]float64, n)
	for i := range durations {
		durations[i] = rand.Float64() * float64(n)
		if durations[i] <= 0 {
			durations[i] = 1e-6
		}
	}
	return durations
}
