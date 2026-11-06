package delivery

import (
	"pr2/internal/insurance"
	"pr2/internal/item"
	"pr2/internal/types"
)

// Базовая реализация доставки
type Delivery struct {
	DefaultCost types.Money
	Insurance   insurance.Insurance
	positions   []item.Position
}

func (d *Delivery) AddItem(t item.Item) {
	d.positions = append(d.positions, item.Position{
		Items: []item.Item{t},
	})
}

func (d *Delivery) AddPosition(p item.Position) {
	d.positions = append(d.positions, p)
}

func (d *Delivery) AddPositions(ps []item.Position) {
	d.positions = append(d.positions, ps...)
}

func (d Delivery) GetPositions() []item.Position {
	return d.positions
}

func (d Delivery) GetPositionsPrice() types.Money {
	sum := types.MoneyZero
	for _, p := range d.positions {
		sum = sum.Add(p.GetPrice())
	}
	return sum
}

func (d Delivery) GetDeliveryPrice() (types.Money, error) {
	return d.DefaultCost, nil
}

func (d Delivery) GetInsurancePrice() types.Money {
	return d.Insurance.Calculate(d.GetPositionsPrice())
}

func (d Delivery) getTotalQuanity() types.Quanity {
	sum := types.QuanityZero
	for _, i := range d.positions {
		sum = sum.Add(i.GetQuanity())
	}
	return sum
}

func (d Delivery) hasItem(f func(item.Item) bool) bool {
	for _, p := range d.positions {
		if p.HasItem(f) {
			return true
		}
	}
	return false
}
