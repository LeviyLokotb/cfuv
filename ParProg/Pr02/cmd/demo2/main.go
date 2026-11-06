// cmd/demo2/main.go
package main

import (
	"fmt"

	"pr2/internal/basket"
	"pr2/internal/delivery"
	"pr2/internal/item"
	"pr2/internal/types"
)

func main() {
	items := []item.Item{
		item.NewItem("Книга", item.TypeBook, types.Quanity(1), types.MoneyFromRub(1000)),
		item.NewItem("Еда", item.TypeFood, types.Quanity(2), types.MoneyFromRub(300)),
	}

	deliveries := []struct {
		name string
		d    delivery.IDelivery
	}{
		{"Simple", makeSimple(items)},
		{"Express", makeExpress(items)},
		{"International", makeInternational(items)},
	}

	for _, entry := range deliveries {
		b := &basket.Basket{}
		b.SetDelivery(entry.d)

		total, err := b.CalculatePrice()
		if err != nil {
			fmt.Printf("%-15s ошибка: %v\n", entry.name, err)
			continue
		}
		fmt.Printf("%-15s итого: %8.2f руб.\n", entry.name, mustRub(total))
	}
}

func makeSimple(items []item.Item) delivery.SimpleDelivery {
	d := delivery.NewDefaultSimpleDelivery()
	for _, it := range items {
		d.Delivery.AddItem(it)
	}
	return d
}

func makeExpress(items []item.Item) delivery.ExpressDelivery {
	d := delivery.NewDefaultExpressDelivery()
	for _, it := range items {
		d.Delivery.AddItem(it)
	}
	return d
}

func makeInternational(items []item.Item) delivery.InternationalDelivery {
	d := delivery.NewDefaultInternationalDelivery()
	for _, it := range items {
		d.Delivery.AddItem(it)
	}
	return d
}

func mustRub(m types.Money) float64 {
	r, _ := m.GetRub()
	return r
}
