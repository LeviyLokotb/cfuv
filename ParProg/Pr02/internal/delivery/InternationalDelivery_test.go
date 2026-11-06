package delivery

import (
	"testing"

	"pr2/internal/item"
	"pr2/internal/types"
)

func TestInternationalDeliveryDefault(t *testing.T) {
	d := NewDefaultInternationalDelivery()
	if len(d.BannedTypes) != 1 || d.BannedTypes[0] != item.TypeEnergyCell {
		t.Errorf("BannedTypes = %v, want [TypeEnergyCell]", d.BannedTypes)
	}
	if d.DefaultCost != types.MoneyFromRub(200) {
		t.Errorf("DefaultCost = %v, want 200_00", d.DefaultCost)
	}
}

func TestInternationalDeliveryAllowed(t *testing.T) {
	d := NewDefaultInternationalDelivery()
	d.Delivery.AddItem(item.NewItem("book", item.TypeBook, 1, types.Money(100)))

	got, err := d.GetDeliveryPrice()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != types.MoneyFromRub(200) {
		t.Errorf("GetDeliveryPrice = %v, want 200_00", got)
	}
}

func TestInternationalDeliveryBanned(t *testing.T) {
	d := NewDefaultInternationalDelivery()
	d.Delivery.AddItem(item.NewItem("battery", item.TypeEnergyCell, 1, types.Money(100)))

	_, err := d.GetDeliveryPrice()
	if err == nil {
		t.Errorf("GetDeliveryPrice: want error for banned TypeEnergyCell")
	}
}

func TestInternationalDeliveryBannedInsidePosition(t *testing.T) {
	d := NewDefaultInternationalDelivery()
	d.Delivery.AddPosition(item.Position{Items: []item.Item{
		item.NewItem("food", item.TypeFood, 1, types.Money(100)),
		item.NewItem("cell", item.TypeEnergyCell, 1, types.Money(100)),
	}})

	_, err := d.GetDeliveryPrice()
	if err == nil {
		t.Errorf("GetDeliveryPrice: want error for banned type inside position")
	}
}

func TestInternationalDeliveryInsurance(t *testing.T) {
	d := NewDefaultInternationalDelivery()
	d.Delivery.AddItem(item.NewItem("book", item.TypeBook, 1, types.MoneyFromRub(10000)))
	// 10000_00 * 0.02 = 200_00; MinCost = 200_00 → 200_00
	got := d.GetInsurancePrice()
	if got != types.MoneyFromRub(200) {
		t.Errorf("GetInsurancePrice = %v, want 200_00", got)
	}
}
