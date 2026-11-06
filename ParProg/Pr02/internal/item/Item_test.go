package item

import (
	"testing"

	"pr2/internal/types"
)

func TestNewItemAndGetters(t *testing.T) {
	it := NewItem("test", TypeFood, types.Quanity(3), types.Money(1000))

	if got := it.GetAmount(); got != types.Quanity(3) {
		t.Errorf("GetAmount = %v, want 3", got)
	}
	if got := it.GetPrice(); got != types.Money(3000) {
		t.Errorf("GetPrice = %v, want 3000", got)
	}
	if got := it.GetType(); got != TypeFood {
		t.Errorf("GetType = %v, want %v", got, TypeFood)
	}
}

func TestItemGetPriceZero(t *testing.T) {
	it := NewItem("t", TypeMisc, types.Quanity(1), types.MoneyZero)
	if got := it.GetPrice(); got != types.MoneyZero {
		t.Errorf("GetPrice = %v, want 0", got)
	}
}
