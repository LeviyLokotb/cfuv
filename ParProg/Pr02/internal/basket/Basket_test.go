package basket

import (
	"testing"

	"pr2/internal/delivery"
	"pr2/internal/item"
	"pr2/internal/types"
)

func TestBasketCalculatePrice(t *testing.T) {
	d := delivery.NewDefaultSimpleDelivery()
	d.Delivery.AddItem(item.NewItem("a", item.TypeBook, 1, types.MoneyFromRub(1000)))

	b := &Basket{}
	b.SetDelivery(d)

	got, err := b.CalculatePrice()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// positions = 1000_00
	// delivery  = 0 (1000_00 > 500_00 → бесплатно)
	// insurance = 1000_00 * 0.01 = 10_00 > min(50_00) → 50_00
	want := types.MoneyFromRub(1000).Add(types.MoneyZero).Add(types.MoneyFromRub(50))
	if got != want {
		t.Errorf("CalculatePrice = %v, want %v", got, want)
	}
}

func TestBasketCalculatePriceWithDeliveryError(t *testing.T) {
	d := delivery.NewDefaultExpressDelivery()
	// 6 штук > 5 → ошибка доставки
	d.Delivery.AddItem(item.NewItem("a", item.TypeMisc, 6, types.MoneyFromRub(100)))

	b := &Basket{}
	b.SetDelivery(d)

	_, err := b.CalculatePrice()
	if err == nil {
		t.Errorf("CalculatePrice: want error")
	}
}

func TestBasketSetDelivery(t *testing.T) {
	b := &Basket{}
	d := delivery.NewDefaultSimpleDelivery()
	b.SetDelivery(d)

	if b.IDelivery == nil {
		t.Errorf("SetDelivery did not set IDelivery")
	}
}
