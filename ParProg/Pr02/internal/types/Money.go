package types

import (
	"errors"
	"math"
)

// Деньги в копейках
type Money int

const (
	// Неверное значение для GetKop
	MoneyInvalid Money = -1
	// Нулевое значение
	MoneyZero Money = 0
)

func MoneyFromRub(rub int) Money {
	return Money(rub * 100)
}

func MoneyFromRubF(rub float64) Money {
	return Money(math.Round(rub * 100.0))
}

func (m Money) GetPercent(p float64) Money {
	n, err := m.GetRub()
	if err != nil {
		return MoneyInvalid
	}
	return MoneyFromRubF(n * p)
}

func (m1 Money) Add(m2 Money) Money {
	p1, err := m1.GetKop()
	if err != nil {
		return MoneyInvalid
	}
	p2, err := m2.GetKop()
	if err != nil {
		return MoneyInvalid
	}
	return Money(p1 + p2)
}

func (m Money) GetKop() (int, error) {
	if m < 0 {
		return 0, errors.New("money must be >= 0")
	}
	return int(m), nil
}

func (m Money) GetRub() (float64, error) {
	k, err := m.GetKop()
	return float64(k) / 100.0, err
}
