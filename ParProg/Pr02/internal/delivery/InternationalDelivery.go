package delivery

import (
	"errors"
	"pr2/internal/insurance"
	"pr2/internal/item"
	"pr2/internal/types"
	"slices"
)

// Доставка стоит DefaultCost.
// Если среди товаров есть товар из категории BannedTypes -- недоступна
type InternationalDelivery struct {
	Delivery
	BannedTypes []item.Type
}

func NewDefaultInternationalDelivery() InternationalDelivery {
	return InternationalDelivery{
		Delivery: Delivery{
			DefaultCost: types.MoneyFromRub(200),
			Insurance: insurance.Insurance{
				MinCost: types.MoneyFromRub(200),
				Percent: 0.02,
			},
		},
		BannedTypes: []item.Type{
			item.TypeEnergyCell,
		},
	}
}

func (d InternationalDelivery) GetDeliveryPrice() (types.Money, error) {
	if d.hasItem(func(i item.Item) bool {
		return slices.Contains(d.BannedTypes, i.GetType())
	}) {
		return types.MoneyInvalid, errors.New("international delivery not allowed")
	}
	return d.DefaultCost, nil
}
