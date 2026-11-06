package geo

import (
	"math"
	"testing"
)

func TestNewRoute(t *testing.T) {
	a := mustCoord(t, 55.7558, 37.6173)
	b := mustCoord(t, 59.9343, 30.3351)

	cases := []struct {
		name    string
		pos     []Coordinate
		wantErr bool
	}{
		{"two points", []Coordinate{a, b}, false},
		{"three points", []Coordinate{a, b, mustCoord(t, 60, 30)}, false},
		{"empty", nil, true},
		{"one point", []Coordinate{a}, true},
		{"identical in a row", []Coordinate{a, a}, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := NewRoute(c.pos, HaversineDistanceCalc{})
			if c.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !c.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestRouteLength(t *testing.T) {
	a := mustCoord(t, 55.7558, 37.6173)
	b := mustCoord(t, 59.9343, 30.3351)
	c := mustCoord(t, 60.0, 30.0)

	r, err := NewRoute([]Coordinate{a, b, c}, HaversineDistanceCalc{})
	if err != nil {
		t.Fatal(err)
	}

	l, err := r.Length()
	if err != nil {
		t.Fatal(err)
	}
	if l <= 0 {
		t.Errorf("length must be positive, got %v", l)
	}

	d1, _ := haversine(a, b)
	d2, _ := haversine(b, c)
	want := d1 + d2
	if math.Abs(l-want) > 1e-6 {
		t.Errorf("Length = %v, want %v", l, want)
	}
}

func TestRouteLengthInvalid(t *testing.T) {
	// Создаём маршрут в обход валидации: одна позиция
	r := Route{
		positions:          []Coordinate{mustCoord(t, 10, 20)},
		DistanceCalculator: HaversineDistanceCalc{},
		valid_state:        state_unknown,
	}
	if _, err := r.Length(); err == nil {
		t.Error("expected error for invalid route")
	}
}

func TestRouteAddPoint(t *testing.T) {
	a := mustCoord(t, 55.7558, 37.6173)
	b := mustCoord(t, 59.9343, 30.3351)
	c := mustCoord(t, 60.0, 30.0)

	r, err := NewRoute([]Coordinate{a, b}, HaversineDistanceCalc{})
	if err != nil {
		t.Fatal(err)
	}

	if err := r.AddPoint(c); err != nil {
		t.Fatal(err)
	}
	if len(r.positions) != 3 {
		t.Errorf("expected 3 positions, got %d", len(r.positions))
	}
}

func TestRouteAddPointIdentical(t *testing.T) {
	a := mustCoord(t, 55.7558, 37.6173)
	b := mustCoord(t, 59.9343, 30.3351)

	r, err := NewRoute([]Coordinate{a, b}, HaversineDistanceCalc{})
	if err != nil {
		t.Fatal(err)
	}
	oldLen := len(r.positions)

	// Добавляем точку, идентичную последней — должно быть ошибкой
	if err := r.AddPoint(b); err == nil {
		t.Error("expected error for identical point")
	}
	if len(r.positions) != oldLen {
		t.Errorf("positions length changed: %d -> %d", oldLen, len(r.positions))
	}
}

// Регрессия: append в AddPoint не должен портить исходный slice при ошибке.
func TestRouteAddPointDoesNotMutateOnError(t *testing.T) {
	a := mustCoord(t, 10, 20)
	b := mustCoord(t, 11, 21)

	positions := make([]Coordinate, 2, 4) // capacity больше len
	positions[0] = a
	positions[1] = b

	r := Route{
		positions:          positions,
		DistanceCalculator: HaversineDistanceCalc{},
		valid_state:        state_unknown,
	}
	if err := r.CheckValid(); err != nil {
		t.Fatal(err)
	}

	// Пытаемся добавить точку, идентичную b — ошибка
	if err := r.AddPoint(b); err == nil {
		t.Error("expected error")
	}

	if len(r.positions) != 2 {
		t.Errorf("positions length changed: %d", len(r.positions))
	}
	if !r.positions[0].Equal(a) || !r.positions[1].Equal(b) {
		t.Error("positions content changed")
	}
}

func TestRouteGetDistanceAfterIdx(t *testing.T) {
	a := mustCoord(t, 55.7558, 37.6173)
	b := mustCoord(t, 59.9343, 30.3351)
	c := mustCoord(t, 60.0, 30.0)

	r, err := NewRoute([]Coordinate{a, b, c}, HaversineDistanceCalc{})
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 2; i++ {
		if _, err := r.GetDistanceAfterIdx(i); err != nil {
			t.Errorf("GetDistanceAfterIdx(%d) error: %v", i, err)
		}
	}
	for _, i := range []int{-1, 2, 100} {
		if _, err := r.GetDistanceAfterIdx(i); err == nil {
			t.Errorf("GetDistanceAfterIdx(%d) expected error", i)
		}
	}
}

// Регрессия: раньше функция всегда возвращала не-nil ошибку.
func TestRouteGetDistanceAfterIdxNoError(t *testing.T) {
	a := mustCoord(t, 55.7558, 37.6173)
	b := mustCoord(t, 59.9343, 30.3351)

	r, _ := NewRoute([]Coordinate{a, b}, HaversineDistanceCalc{})
	d, err := r.GetDistanceAfterIdx(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d <= 0 {
		t.Errorf("distance must be positive, got %v", d)
	}
}

func TestRouteGetDistanceAfter(t *testing.T) {
	a := mustCoord(t, 55.7558, 37.6173)
	b := mustCoord(t, 59.9343, 30.3351)
	c := mustCoord(t, 60.0, 30.0)

	r, _ := NewRoute([]Coordinate{a, b, c}, HaversineDistanceCalc{})

	if _, err := r.GetDistanceAfter(a); err != nil {
		t.Errorf("GetDistanceAfter(a) error: %v", err)
	}
	if _, err := r.GetDistanceAfter(b); err != nil {
		t.Errorf("GetDistanceAfter(b) error: %v", err)
	}
	// Последняя точка не имеет «следующей»
	if _, err := r.GetDistanceAfter(c); err == nil {
		t.Error("GetDistanceAfter(last) expected error")
	}
	// Точки, которой нет в маршруте
	if _, err := r.GetDistanceAfter(mustCoord(t, 0, 0)); err == nil {
		t.Error("GetDistanceAfter(unknown) expected error")
	}
}

func TestRouteCheckValidCaching(t *testing.T) {
	a := mustCoord(t, 10, 20)
	b := mustCoord(t, 11, 21)

	r, err := NewRoute([]Coordinate{a, b}, HaversineDistanceCalc{})
	if err != nil {
		t.Fatal(err)
	}
	if r.valid_state != state_valid {
		t.Errorf("expected state_valid, got %v", r.valid_state)
	}

	if err := r.CheckValid(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if r.valid_state != state_valid {
		t.Errorf("state changed after second CheckValid: %v", r.valid_state)
	}
}
