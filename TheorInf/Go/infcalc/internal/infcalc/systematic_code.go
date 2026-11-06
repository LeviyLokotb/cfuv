package infcalc

import (
	"theorinf/internal/matrix"
)

/*
p. Минимальная длина систематического кода
для обнаружения и исправления однократных ошибок

k -- длина исходного сообщения
*/
func MinSyscodeLen(k int) int {
	p := 1
	for p*p-p-1 > k {
		p++
	}
	return p
}

/*
n. Длина сообщения с систематическим кодом

k -- длина исходного сообщения

p -- длина систематического кода
*/
func MessageWithSyscodeLen(k int, p int) int {
	return k + p
}

/*
U_k. Левая подматрица k x k производящей матрицы P_n,k.
Отвечает за информационные разряды.

Заполняется как единичная матрица
*/
func GenerateSyscodeInfoSubmx(k int) [][]float64 {
	return matrix.E(k)
}

/*
H_p. Правая подматрица k x p производящей матрицы P_n,k.
Отвечает за проверочные разряды.

k -- длина исходного сообщения

p -- длина систематического кода

Заполняется в соответствии с правилами:

- строки ненулевые

- строки попарно различаются (разница хотя бы 1 бит)
*/
func GenerateSyscodeCodeSubmx(k int, p int) {
	H := make([][]float64, k)
	for i := range H {
		row := make([]float64, p)
		for {
			// 
			if !matrix.Contains(H, row) {
				break
			}
		}
		H[i] = row
	}
}
