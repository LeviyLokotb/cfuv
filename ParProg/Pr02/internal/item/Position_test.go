package item

import (
	"testing"

	"pr2/internal/types"
)

func TestPositionGetPrice(t *testing.T) {
	p := Position{Items: []Item{
		NewItem("a", TypeFood, types.Quanity(2), types.Money(100)),
		NewItem("b", TypeBook, types.Quanity(1), types.Money(500)),
	}}
	if got := p.GetPrice(); got != types.Money(700) {
		t.Errorf("GetPrice = %v, want 700", got)
	}
}

func TestPositionGetPriceEmpty(t *testing.T) {
	if got := (Position{}).GetPrice(); got != types.MoneyZero {
		t.Errorf("empty GetPrice = %v, want 0", got)
	}
}

func TestPositionGetQuanity(t *testing.T) {
	p := Position{Items: []Item{
		NewItem("a", TypeFood, types.Quanity(2), types.Money(100)),
		NewItem("b", TypeBook, types.Quanity(3), types.Money(500)),
	}}
	if got := p.GetQuanity(); got != types.Quanity(5) {
		t.Errorf("GetQuanity = %v, want 5", got)
	}
}

func TestPositionHasItem(t *testing.T) {
	p := Position{Items: []Item{
		NewItem("a", TypeFood, types.Quanity(1), types.Money(100)),
		NewItem("b", TypeBook, types.Quanity(1), types.Money(500)),
	}}

	if !p.HasItem(func(i Item) bool { return i.GetType() == TypeBook }) {
		t.Errorf("HasItem(TypeBook) = false, want true")
	}
	if p.HasItem(func(i Item) bool { return i.GetType() == TypeMeds }) {
		t.Errorf("HasItem(TypeMeds) = true, want false")
	}
}

func TestPositionHasItemEmpty(t *testing.T) {
	p := Position{}
	if p.HasItem(func(i Item) bool { return true }) {
		t.Errorf("empty HasItem = true, want false")
	}
}
