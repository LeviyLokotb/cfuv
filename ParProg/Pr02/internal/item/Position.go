package item

import "pr2/internal/types"

type Position struct {
	Items []Item
}

func (p Position) GetPrice() types.Money {
	sum := types.MoneyZero
	for _, i := range p.Items {
		sum = sum.Add(i.GetPrice())
	}
	return sum
}

func (p Position) GetQuanity() types.Quanity {
	sum := types.QuanityZero
	for _, i := range p.Items {
		sum = sum.Add(i.quanity)
	}
	return sum
}

func (p Position) HasItem(f func(Item) bool) bool {
	for _, i := range p.Items {
		if f(i) {
			return true
		}
	}
	return false
}
