package infcalc

import "math"

/*
H(Y/X). Условная энтропия -- остаточная неопределённость,
энтропия шума

condProbs -- матрица условных вероятностей P(Y/X)

jointProbs -- матрица совместных вероятностей P(X,Y)
*/
func ConditionalEntropy(condProbs [][]float64, jointProbs [][]float64, n int) float64 {
	var sum float64
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if condProbs[i][j] == 0 {
				continue
			}
			sum += jointProbs[i][j] * math.Log2(condProbs[i][j])
		}
	}
	return -sum
}

/*
I(X,Y). Среднее количество информации, получаемое при неполной достоверности сообщений
(Взаимная информация, сколько дошло от входа к выходу)

entropyOut -- энтропия на выходе H(Y)

conditionalEntropy -- условная энтропия H(Y/X)
*/
func MutualInformation(entropyOut float64, conditionalEntropy float64) float64 {
	return entropyOut - conditionalEntropy
}

/*
τ̅. Средняя длительность символа

probsIn -- массив вероятностей появления сообщений на входе P(X)

durations -- массив длительностей символов τ
*/
func MeanDuration(probsIn []float64, durations []float64, n int) float64 {
	var sum float64
	for i := 0; i < n; i++ {
		sum += probsIn[i] * durations[i]
	}
	return sum
}

/*
C. Пропускная способность дискретного канала без помех

meanDur -- средняя длительность символа τ̅
*/
func CapacityNoNoise(n int, meanDur float64) float64 {
	return CapacityWithNoise(n, 0, meanDur)
}

/*
C. Пропускная способность дискретного канала с помехами

entropyNoise -- энтропия шума H(Y/X)

meanDur -- средняя длительность символа τ̅
*/
func CapacityWithNoise(n int, entropyNoise float64, meanDur float64) float64 {
	if meanDur <= 0 {
		return 0
	}
	return (math.Log2(float64(n)) - entropyNoise) / meanDur
}

/*
I̅(Y). Скорость передачи для дискретного канала без помех

entropyIn -- энтропия на входе H(X) либо на выходе H(Y) (совпадают)

meanDur -- средняя длительность символа τ̅
*/
func SpeedNoNoise(entropyIn float64, meanDur float64) float64 {
	return SpeedWithNoise(entropyIn, 0, meanDur)
}

/*
I̅(Z,Y). Скорость передачи для дискретного канала с помехами

entropyOut -- энтропия на выходе H(Y)

entropyNoise -- энтропия шума H(Y/X)

meanDur -- средняя длительность символа τ̅
*/
func SpeedWithNoise(entropyOut float64, entropyNoise float64, meanDur float64) float64 {
	if meanDur <= 0 {
		return 0
	}
	return (entropyOut - entropyNoise) / meanDur
}
