package insurance

import "pr2/internal/types"

type Insurance struct {
	MinCost types.Money
	Percent float64
}

func (s Insurance) Calculate(price types.Money) types.Money {
	p := price.GetPercent(s.Percent)
	// Обрабатываем ошибочный результат
	if _, err := p.GetKop(); err != nil {
		return p
	}

	if p > s.MinCost {
		return p
	}
	return s.MinCost
}
