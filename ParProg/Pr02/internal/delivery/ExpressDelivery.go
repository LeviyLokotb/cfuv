package delivery

import (
	"errors"
	"pr2/internal/insurance"
	"pr2/internal/types"
)

// Доставка стоит DefaultCost. Если количество товаров больше EdgeQuanity -- недоступна.
// Страховка -- Percent доля от суммы товаров, минимум MinCost.
// Страховка обязательна при сумме больше InsuranceEdgePrice
type ExpressDelivery struct {
	Delivery
	InsuranceEdgePrice types.Money
	EdgeQuanity        types.Quanity
	InsuranceEnabled   bool
}

func NewDefaultExpressDelivery() ExpressDelivery {
	return ExpressDelivery{
		Delivery: Delivery{
			DefaultCost: types.MoneyFromRub(500),
			Insurance: insurance.Insurance{
				MinCost: types.MoneyZero,
				Percent: 0.015,
			},
		},
		InsuranceEdgePrice: types.MoneyFromRub(10_000),
		EdgeQuanity:        types.Quanity(5),
		InsuranceEnabled:   true,
	}
}

func (d ExpressDelivery) GetInsurancePrice() types.Money {
	if d.InsuranceEnabled || d.GetPositionsPrice() > d.InsuranceEdgePrice {
		return d.Insurance.Calculate(d.GetPositionsPrice())
	}
	return types.MoneyZero
}

func (d ExpressDelivery) GetDeliveryPrice() (types.Money, error) {
	q := d.getTotalQuanity()
	if q > d.EdgeQuanity {
		return types.MoneyInvalid, errors.New("express delivery not allowed")
	}
	return d.DefaultCost, nil
}
