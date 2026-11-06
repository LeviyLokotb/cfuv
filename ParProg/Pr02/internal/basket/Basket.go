package basket

import (
	"fmt"
	"pr2/internal/delivery"
	"pr2/internal/types"
)

type Basket struct {
	delivery.IDelivery
}

func (b *Basket) SetDelivery(d delivery.IDelivery) {
	b.IDelivery = d
}

func (b Basket) CalculatePrice() (types.Money, error) {
	sumdel, err := b.GetDeliveryPrice()
	if err != nil {
		return types.MoneyInvalid, fmt.Errorf("delivery failed: %w", err)
	}
	sumpos := b.GetPositionsPrice()
	sumins := b.GetInsurancePrice()

	return sumpos.Add(sumdel).Add(sumins), nil
}
