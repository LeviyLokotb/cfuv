// cmd/demo3/main.go
package main

import (
	"fmt"

	"pr2/internal/basket"
	"pr2/internal/delivery"
	"pr2/internal/item"
	"pr2/internal/types"
)

func main() {
	// Express: >5 единиц → недоступна
	d := delivery.NewDefaultExpressDelivery()
	d.Delivery.AddItem(item.NewItem("Болт", item.TypeMisc, types.Quanity(6), types.MoneyFromRub(10)))

	b := &basket.Basket{}
	b.SetDelivery(d)

	_, err := b.CalculatePrice()
	if err != nil {
		fmt.Println("Ожидаемая ошибка:", err)
		return
	}
	fmt.Println("Ошибки не было — что-то не так")
}
