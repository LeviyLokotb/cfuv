package delivery

import (
	"pr2/internal/item"
	"pr2/internal/types"
)

type IDelivery interface {
	GetPositions() []item.Position
	GetPositionsPrice() types.Money
	GetInsurancePrice() types.Money
	GetDeliveryPrice() (types.Money, error)
}
