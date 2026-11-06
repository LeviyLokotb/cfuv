package delivery

import (
	"pr2/internal/insurance"
	"pr2/internal/types"
)

// Доставка стоит DefaultCost. Если товары стоят больше EdgeCost -- бесплатно.
// Страховка -- Percent доля от суммы товаров, минимум MinCost
// Страховка включается по желанию
type SimpleDelivery struct {
	Delivery
	EdgeCost         types.Money
	InsuranceEnabled bool
}

func NewDefaultSimpleDelivery() SimpleDelivery {
	return SimpleDelivery{
		Delivery: Delivery{
			DefaultCost: types.MoneyFromRub(200),
			Insurance: insurance.Insurance{
				MinCost: types.MoneyFromRub(50),
				Percent: 0.01,
			},
		},
		EdgeCost:         types.MoneyFromRub(500),
		InsuranceEnabled: true,
	}
}

func (d SimpleDelivery) GetInsurancePrice() types.Money {
	if d.InsuranceEnabled {
		return d.Insurance.Calculate(d.GetPositionsPrice())
	}
	return types.MoneyZero
}

func (d SimpleDelivery) GetDeliveryPrice() (types.Money, error) {
	sum := d.GetPositionsPrice()
	if sum > d.EdgeCost {
		return types.MoneyZero, nil
	}
	return d.DefaultCost, nil
}
