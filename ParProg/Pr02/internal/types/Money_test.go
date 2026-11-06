package types

import (
	"testing"
)

func TestMoneyFromRub(t *testing.T) {
	if got := MoneyFromRub(5); got != Money(500) {
		t.Errorf("MoneyFromRub(5) = %d, want 500", got)
	}
	if got := MoneyFromRub(0); got != MoneyZero {
		t.Errorf("MoneyFromRub(0) = %d, want 0", got)
	}
}

func TestMoneyFromRubF(t *testing.T) {
	cases := []struct {
		in   float64
		want Money
	}{
		{1.0, 100},
		{1.005, 100}, // round-half-to-even? math.Round → 101? проверим
		{1.999, 200},
		{0.0, 0},
		{2.5, 250},
	}
	for _, c := range cases {
		got := MoneyFromRubF(c.in)
		if got != c.want {
			t.Errorf("MoneyFromRubF(%v) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestMoneyAdd(t *testing.T) {
	if got := Money(100).Add(Money(50)); got != Money(150) {
		t.Errorf("100+50 = %d, want 150", got)
	}
	if got := MoneyZero.Add(MoneyZero); got != MoneyZero {
		t.Errorf("0+0 = %d, want 0", got)
	}
	// отрицательные
	if got := MoneyInvalid.Add(Money(10)); got != MoneyInvalid {
		t.Errorf("invalid + 10 = %d, want %d", got, MoneyInvalid)
	}
	if got := Money(10).Add(MoneyInvalid); got != MoneyInvalid {
		t.Errorf("10 + invalid = %d, want %d", got, MoneyInvalid)
	}
}

func TestMoneyGetKop(t *testing.T) {
	if v, err := Money(150).GetKop(); err != nil || v != 150 {
		t.Errorf("GetKop() = %d, %v; want 150, nil", v, err)
	}
	if _, err := MoneyInvalid.GetKop(); err == nil {
		t.Errorf("GetKop() on invalid: want error, got nil")
	}
	if _, err := Money(-5).GetKop(); err == nil {
		t.Errorf("GetKop() on negative: want error, got nil")
	}
}

func TestMoneyGetRub(t *testing.T) {
	v, err := Money(150).GetRub()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != 1.5 {
		t.Errorf("GetRub() = %v, want 1.5", v)
	}
}

func TestMoneyGetPercent(t *testing.T) {
	if got := Money(1000).GetPercent(0.1); got != Money(100) {
		t.Errorf("1000 * 10%% = %d, want 100", got)
	}
	if got := MoneyInvalid.GetPercent(0.1); got != MoneyInvalid {
		t.Errorf("invalid.GetPercent = %d, want invalid", got)
	}
}
