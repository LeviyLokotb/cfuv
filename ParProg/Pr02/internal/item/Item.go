package item

import "pr2/internal/types"

type Item struct {
	title   string
	_type   Type
	quanity types.Quanity
	cost    types.Money
}

func NewItem(
	title string,
	t Type,
	q types.Quanity,
	cost types.Money,
) Item {
	return Item{
		title:   title,
		_type:   t,
		quanity: q,
		cost:    cost,
	}
}

func (i Item) GetAmount() types.Quanity {
	return i.quanity
}

func (i Item) GetPrice() types.Money {
	return i.cost * types.Money(i.quanity)
}

func (i Item) GetType() Type {
	return i._type
}
