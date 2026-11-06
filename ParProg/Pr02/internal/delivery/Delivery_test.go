package delivery

import (
	"testing"

	"pr2/internal/insurance"
	"pr2/internal/item"
	"pr2/internal/types"
)

func newTestItem(t string, cost types.Money, q types.Quanity) item.Item {
	return item.NewItem(t, item.TypeMisc, q, cost)
}

func TestDeliveryAddItem(t *testing.T) {
	d := &Delivery{DefaultCost: types.Money(100)}
	d.AddItem(newTestItem("a", 100, 1))
	d.AddItem(newTestItem("b", 200, 1))

	if got := len(d.GetPositions()); got != 2 {
		t.Errorf("positions = %d, want 2", got)
	}
}

func TestDeliveryAddPosition(t *testing.T) {
	d := &Delivery{}
	p := item.Position{Items: []item.Item{
		newTestItem("a", 100, 1),
		newTestItem("b", 200, 1),
	}}
	d.AddPosition(p)

	if got := len(d.GetPositions()); got != 1 {
		t.Errorf("positions = %d, want 1", got)
	}
}

func TestDeliveryAddPositions(t *testing.T) {
	d := &Delivery{}
	ps := []item.Position{
		{Items: []item.Item{newTestItem("a", 100, 1)}},
		{Items: []item.Item{newTestItem("b", 200, 1)}},
	}
	d.AddPositions(ps)

	if got := len(d.GetPositions()); got != 2 {
		t.Errorf("positions = %d, want 2", got)
	}
}

func TestDeliveryGetPositionsPrice(t *testing.T) {
	d := &Delivery{}
	d.AddItem(newTestItem("a", 100, 2)) // 200
	d.AddItem(newTestItem("b", 500, 1)) // 500

	if got := d.GetPositionsPrice(); got != types.Money(700) {
		t.Errorf("GetPositionsPrice = %v, want 700", got)
	}
}

func TestDeliveryGetDeliveryPrice(t *testing.T) {
	d := Delivery{DefaultCost: types.MoneyFromRub(500)}
	got, err := d.GetDeliveryPrice()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != types.MoneyFromRub(500) {
		t.Errorf("GetDeliveryPrice = %v, want 500_00", got)
	}
}

func TestDeliveryGetInsurancePrice(t *testing.T) {
	d := &Delivery{
		Insurance: insurance.Insurance{MinCost: types.Money(50), Percent: 0.01},
	}
	d.AddItem(newTestItem("a", 1000, 1)) // 1000 * 0.01 = 10 < 50 → 50

	if got := d.GetInsurancePrice(); got != types.Money(50) {
		t.Errorf("GetInsurancePrice = %v, want 50", got)
	}
}
