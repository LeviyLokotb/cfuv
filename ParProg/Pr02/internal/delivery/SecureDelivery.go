package delivery

import (
	"errors"
	"pr2/internal/insurance"
	"pr2/internal/types"
)

// Доставка стоит DefaultCost + AdditionalPercent от стоимости товаров.
// Если стоимость товаров больше MaxPrice -- недоступна.
// Если нет страховки -- недоступна.
type SecureDelivery struct {
	Delivery
	MaxPrice          types.Money
	AdditionalPercent float64
	InsuranceEnabled  bool
}

func NewDefaultSecureDelivery() SecureDelivery {
	return SecureDelivery{
		Delivery: Delivery{
			DefaultCost: types.MoneyFromRub(1000),
			Insurance: insurance.Insurance{
				MinCost: types.MoneyFromRub(50),
				Percent: 0.01,
			},
		},
		InsuranceEnabled:  true,
		MaxPrice:          types.MoneyFromRub(100_000),
		AdditionalPercent: 0.005,
	}
}

func (d SecureDelivery) GetInsurancePrice() types.Money {
	if d.InsuranceEnabled {
		return d.Insurance.Calculate(d.GetPositionsPrice())
	}
	return types.MoneyZero
}

func (d SecureDelivery) GetDeliveryPrice() (types.Money, error) {
	// Доставка недоступна если нет страхования
	if !d.InsuranceEnabled {
		return types.MoneyInvalid, errors.New("secure delivery is not allowed")
	}
	// Стоимость товаров, не больше MaxPrice
	price := d.GetPositionsPrice()
	if price > d.MaxPrice {
		return types.MoneyInvalid, errors.New("secure delivery is not allowed")
	}
	// Дополнительный процент от стоимости
	proc := price.GetPercent(d.AdditionalPercent)
	return d.DefaultCost + proc, nil
}
