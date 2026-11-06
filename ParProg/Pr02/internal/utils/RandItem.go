package utils

import (
	"fmt"
	"math/rand"
	"pr2/internal/item"
	"pr2/internal/types"
	"time"
)

var id int64 = 0

func RandItem() item.Item {
	title := fmt.Sprintf("item#%d:%d", time.Now().UnixNano(), id)
	id++

	t := rand.Intn(int(item.TypeOther) + 1)

	q := rand.Intn(11)
	for q%10 == 0 {
		q += rand.Intn(11)
	}

	cost := rand.Float64()*9900 + 100

	return item.NewItem(
		title,
		item.Type(t),
		types.Quanity(q),
		types.Money(cost),
	)
}
