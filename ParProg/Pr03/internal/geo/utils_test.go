package geo

import (
	"math"
	"testing"
)

func mustCoord(t *testing.T, lat, lon float64) Coordinate {
	t.Helper()
	c, err := NewCoordinateFromDeg(lat, lon)
	if err != nil {
		t.Fatalf("NewCoordinateFromDeg(%v,%v): %v", lat, lon, err)
	}
	return c
}

func TestHaversineKnownDistances(t *testing.T) {
	cases := []struct {
		name       string
		lat1, lon1 float64
		lat2, lon2 float64
		wantM      float64
		tolM       float64
	}{
		{
			name: "Moscow - SPb",
			lat1: 55.7558, lon1: 37.6173,
			lat2: 59.9343, lon2: 30.3351,
			wantM: 635_000, tolM: 10_000,
		},
		{
			name: "same point",
			lat1: 10, lon1: 20,
			lat2: 10, lon2: 20,
			wantM: 0, tolM: 1e-6,
		},
		{
			name: "1 degree latitude",
			lat1: 0, lon1: 0,
			lat2: 1, lon2: 0,
			wantM: earth_radius_m * math.Pi / 180, tolM: 100,
		},
		{
			name: "equator quarter",
			lat1: 0, lon1: 0,
			lat2: 0, lon2: 90,
			wantM: earth_radius_m * math.Pi / 2, tolM: 1000,
		},
		{
			name: "antipodal",
			lat1: 0, lon1: 0,
			lat2: 0, lon2: 180,
			wantM: math.Pi * earth_radius_m, tolM: 1,
		},
		{
			name: "north pole to south pole",
			lat1: 90, lon1: 0,
			lat2: -90, lon2: 0,
			wantM: math.Pi * earth_radius_m, tolM: 1,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			a := mustCoord(t, c.lat1, c.lon1)
			b := mustCoord(t, c.lat2, c.lon2)
			d, err := haversine(a, b)
			if err != nil {
				t.Fatal(err)
			}
			if math.Abs(d-c.wantM) > c.tolM {
				t.Errorf("got %v, want %v (±%v)", d, c.wantM, c.tolM)
			}
		})
	}
}

func TestHaversineSymmetry(t *testing.T) {
	a := mustCoord(t, 10, 20)
	b := mustCoord(t, 30, 40)
	d1, _ := haversine(a, b)
	d2, _ := haversine(b, a)
	if math.Abs(d1-d2) > 1e-6 {
		t.Errorf("not symmetric: %v vs %v", d1, d2)
	}
}

// Проверка клипа na: пары, где na может выйти за [0,1] из-за погрешности.
func TestHaversineNoNaN(t *testing.T) {
	pairs := [][4]float64{
		{90, 0, 90, 0},
		{-90, 0, -90, 0},
		{90, 0, -90, 180},
		{0, 0, 0, 180},
		{89.9999999, 0, -89.9999999, 180},
	}
	for _, p := range pairs {
		a := mustCoord(t, p[0], p[1])
		b := mustCoord(t, p[2], p[3])
		d, err := haversine(a, b)
		if err != nil {
			t.Fatalf("haversine(%v): %v", p, err)
		}
		if math.IsNaN(d) || math.IsInf(d, 0) {
			t.Errorf("haversine(%v) = %v (NaN/Inf)", p, d)
		}
	}
}

func TestSphereCosTheoremKnownDistances(t *testing.T) {
	cases := []struct {
		name       string
		lat1, lon1 float64
		lat2, lon2 float64
		wantM      float64
		tolM       float64
	}{
		{
			name: "Moscow - SPb",
			lat1: 55.7558, lon1: 37.6173,
			lat2: 59.9343, lon2: 30.3351,
			wantM: 635_000, tolM: 10_000,
		},
		{
			name: "same point",
			lat1: 10, lon1: 20,
			lat2: 10, lon2: 20,
			wantM: 0, tolM: 1e-3,
		},
		{
			name: "antipodal",
			lat1: 0, lon1: 0,
			lat2: 0, lon2: 180,
			wantM: math.Pi * earth_radius_m, tolM: 1,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			a := mustCoord(t, c.lat1, c.lon1)
			b := mustCoord(t, c.lat2, c.lon2)
			d, err := sphereCosTheorem(a, b)
			if err != nil {
				t.Fatal(err)
			}
			if math.Abs(d-c.wantM) > c.tolM {
				t.Errorf("got %v, want %v (±%v)", d, c.wantM, c.tolM)
			}
		})
	}
}

// Проверка клипа аргумента Acos: |s+c| может выйти за [-1,1] из-за погрешности.
func TestSphereCosTheoremNoNaN(t *testing.T) {
	pairs := [][4]float64{
		{90, 0, 90, 0},
		{-90, 0, -90, 0},
		{90, 0, -90, 180},
		{0, 0, 0, 180},
	}
	for _, p := range pairs {
		a := mustCoord(t, p[0], p[1])
		b := mustCoord(t, p[2], p[3])
		d, err := sphereCosTheorem(a, b)
		if err != nil {
			t.Fatalf("sphereCosTheorem(%v): %v", p, err)
		}
		if math.IsNaN(d) || math.IsInf(d, 0) {
			t.Errorf("sphereCosTheorem(%v) = %v (NaN/Inf)", p, d)
		}
	}
}
