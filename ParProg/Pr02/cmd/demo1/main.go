// cmd/demo1/main.go
package main

import (
	"fmt"

	"pr2/internal/basket"
	"pr2/internal/delivery"
	"pr2/internal/item"
	"pr2/internal/types"
)

func main() {
	d := delivery.NewDefaultSimpleDelivery()
	d.Delivery.AddItem(item.NewItem("Книга", item.TypeBook, types.Quanity(1), types.MoneyFromRub(1000)))

	b := &basket.Basket{}
	b.SetDelivery(d)

	total, err := b.CalculatePrice()
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	fmt.Printf("Итого: %v коп. (%.2f руб.)\n", total, mustRub(total))
}

func mustRub(m types.Money) float64 {
	r, _ := m.GetRub()
	return r
}
