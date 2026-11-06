package delivery

import (
	"testing"

	"pr2/internal/types"
)

func TestExpressDeliveryDefault(t *testing.T) {
	d := NewDefaultExpressDelivery()
	if d.EdgeQuanity != types.Quanity(5) {
		t.Errorf("EdgeQuanity = %v, want 5", d.EdgeQuanity)
	}
	if !d.InsuranceEnabled {
		t.Errorf("InsuranceEnabled = false, want true")
	}
}

func TestExpressDeliveryWithinQuanity(t *testing.T) {
	d := NewDefaultExpressDelivery()
	for i := 0; i < 5; i++ {
		d.Delivery.AddItem(newTestItem("a", types.MoneyFromRub(100), 1))
	}

	got, err := d.GetDeliveryPrice()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != types.MoneyFromRub(500) {
		t.Errorf("GetDeliveryPrice = %v, want 500_00", got)
	}
}

func TestExpressDeliveryOverQuanity(t *testing.T) {
	d := NewDefaultExpressDelivery()
	d.Delivery.AddItem(newTestItem("a", types.MoneyFromRub(100), 6))

	_, err := d.GetDeliveryPrice()
	if err == nil {
		t.Errorf("GetDeliveryPrice: want error for quantity 6 > 5")
	}
}

func TestExpressDeliveryInsuranceEnabled(t *testing.T) {
	d := NewDefaultExpressDelivery()
	d.Delivery.AddItem(newTestItem("a", types.MoneyFromRub(10000), 1))
	// 10000_00 * 0.015 = 150_00 > MinCost(0) → 150_00
	got := d.GetInsurancePrice()
	if got != types.MoneyFromRub(150) {
		t.Errorf("GetInsurancePrice = %v, want 150_00", got)
	}
}
