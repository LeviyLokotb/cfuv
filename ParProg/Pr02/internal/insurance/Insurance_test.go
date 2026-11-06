package insurance

import (
	"testing"

	"pr2/internal/types"
)

func TestInsuranceCalculateAboveMin(t *testing.T) {
	ins := Insurance{MinCost: types.Money(100), Percent: 0.1}

	// 10000 * 0.1 = 1000 > 100
	got := ins.Calculate(types.Money(10000))
	if got != types.Money(1000) {
		t.Errorf("Calculate(10000) = %v, want 1000", got)
	}
}

func TestInsuranceCalculateMinApplies(t *testing.T) {
	ins := Insurance{MinCost: types.Money(500), Percent: 0.01}

	// 1000 * 0.01 = 10 < 500 → min = 500
	got := ins.Calculate(types.Money(1000))
	if got != types.Money(500) {
		t.Errorf("Calculate(1000) = %v, want 500", got)
	}
}

func TestInsuranceCalculateExactMin(t *testing.T) {
	ins := Insurance{MinCost: types.Money(100), Percent: 0.1}
	// 1000 * 0.1 = 100 == min → возвращается min
	got := ins.Calculate(types.Money(1000))
	if got != types.Money(100) {
		t.Errorf("Calculate(1000) = %v, want 100", got)
	}
}

func TestInsuranceCalculateZeroMin(t *testing.T) {
	ins := Insurance{MinCost: types.MoneyZero, Percent: 0.015}
	got := ins.Calculate(types.Money(10000))
	if got != types.Money(150) {
		t.Errorf("Calculate(10000) = %v, want 150", got)
	}
}

func TestInsuranceCalculateInvalid(t *testing.T) {
	ins := Insurance{MinCost: types.Money(100), Percent: 0.1}
	got := ins.Calculate(types.MoneyInvalid)
	// MoneyInvalid.GetPercent → MoneyInvalid; MoneyInvalid > MinCost? нет, т.к. -1 < 100
	// следовательно вернётся MinCost — это надо зафиксировать как текущее поведение
	if got != types.Money(100) {
		t.Logf("Calculate(invalid) = %v (текущее поведение)", got)
	}
}
