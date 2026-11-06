package delivery

import (
	"testing"

	"pr2/internal/types"
)

func TestSimpleDeliveryDefault(t *testing.T) {
	d := NewDefaultSimpleDelivery()

	if d.DefaultCost != types.MoneyFromRub(200) {
		t.Errorf("DefaultCost = %v, want %v", d.DefaultCost, types.MoneyFromRub(200))
	}
	if d.EdgeCost != types.MoneyFromRub(500) {
		t.Errorf("EdgeCost = %v, want %v", d.EdgeCost, types.MoneyFromRub(500))
	}
	if !d.InsuranceEnabled {
		t.Errorf("InsuranceEnabled = false, want true")
	}
	if d.Insurance.MinCost != types.MoneyFromRub(50) {
		t.Errorf("Insurance.MinCost = %v, want %v", d.Insurance.MinCost, types.MoneyFromRub(50))
	}
	if d.Insurance.Percent != 0.01 {
		t.Errorf("Insurance.Percent = %v, want 0.01", d.Insurance.Percent)
	}
}

func TestSimpleDeliveryPriceBelowEdge(t *testing.T) {
	d := NewDefaultSimpleDelivery()
	// сумма = 100_00 < EdgeCost(500_00) → платим DefaultCost
	d.Delivery.AddItem(newTestItem("a", types.MoneyFromRub(100), 1))

	got, err := d.GetDeliveryPrice()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != types.MoneyFromRub(200) {
		t.Errorf("GetDeliveryPrice = %v, want %v", got, types.MoneyFromRub(200))
	}
}

func TestSimpleDeliveryPriceAboveEdgeIsFree(t *testing.T) {
	d := NewDefaultSimpleDelivery()
	// сумма = 1000_00 > EdgeCost(500_00) → бесплатно
	d.Delivery.AddItem(newTestItem("a", types.MoneyFromRub(1000), 1))

	got, err := d.GetDeliveryPrice()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != types.MoneyZero {
		t.Errorf("GetDeliveryPrice = %v, want 0", got)
	}
}

func TestSimpleDeliveryPriceExactlyEdge(t *testing.T) {
	d := NewDefaultSimpleDelivery()
	// сумма == EdgeCost(500_00): строгое > не срабатывает → платим DefaultCost
	d.Delivery.AddItem(newTestItem("a", types.MoneyFromRub(500), 1))

	got, err := d.GetDeliveryPrice()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != types.MoneyFromRub(200) {
		t.Errorf("GetDeliveryPrice = %v, want %v (на границе платим)",
			got, types.MoneyFromRub(200))
	}
}

func TestSimpleDeliveryPriceEmptyPositions(t *testing.T) {
	d := NewDefaultSimpleDelivery()
	// нет позиций: сумма = 0 < EdgeCost → платим DefaultCost
	got, err := d.GetDeliveryPrice()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != types.MoneyFromRub(200) {
		t.Errorf("GetDeliveryPrice (empty) = %v, want %v",
			got, types.MoneyFromRub(200))
	}
}

func TestSimpleDeliveryInsuranceEnabled(t *testing.T) {
	d := NewDefaultSimpleDelivery()
	d.Delivery.AddItem(newTestItem("a", types.MoneyFromRub(10000), 1))

	// 10000_00 * 0.01 = 100_00 > MinCost(50_00) → 100_00
	got := d.GetInsurancePrice()
	if got != types.MoneyFromRub(100) {
		t.Errorf("GetInsurancePrice = %v, want %v", got, types.MoneyFromRub(100))
	}
}

func TestSimpleDeliveryInsuranceMinApplies(t *testing.T) {
	d := NewDefaultSimpleDelivery()
	d.Delivery.AddItem(newTestItem("a", types.MoneyFromRub(1000), 1))

	// 1000_00 * 0.01 = 10_00 < MinCost(50_00) → 50_00
	got := d.GetInsurancePrice()
	if got != types.MoneyFromRub(50) {
		t.Errorf("GetInsurancePrice = %v, want %v (min)", got, types.MoneyFromRub(50))
	}
}

func TestSimpleDeliveryInsuranceDisabled(t *testing.T) {
	d := NewDefaultSimpleDelivery()
	d.InsuranceEnabled = false
	d.Delivery.AddItem(newTestItem("a", types.MoneyFromRub(10000), 1))

	got := d.GetInsurancePrice()
	if got != types.MoneyZero {
		t.Errorf("GetInsurancePrice (disabled) = %v, want 0", got)
	}
}

func TestSimpleDeliveryInsuranceEmptyPositions(t *testing.T) {
	d := NewDefaultSimpleDelivery()
	// нет позиций: 0 * 0.01 = 0 < MinCost(50_00) → 50_00
	got := d.GetInsurancePrice()
	if got != types.MoneyFromRub(50) {
		t.Errorf("GetInsurancePrice (empty) = %v, want %v (min)",
			got, types.MoneyFromRub(50))
	}
}

// Проверка, что SimpleDelivery удовлетворяет IDelivery на этапе компиляции.
func TestSimpleDeliveryImplementsIDelivery(t *testing.T) {
	var _ IDelivery = SimpleDelivery{}
	var _ IDelivery = NewDefaultSimpleDelivery()
}
