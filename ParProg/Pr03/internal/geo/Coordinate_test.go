package geo

import (
	"math"
	"testing"
)

func TestNewCoordinateFromDeg(t *testing.T) {
	cases := []struct {
		name    string
		lat     float64
		lon     float64
		wantErr bool
	}{
		{"typical", 55.7558, 37.6173, false},
		{"zero", 0, 0, false},
		{"extremes", 90, 180, false},
		{"bad lat", 91, 0, true},
		{"bad lon", 0, 181, true},
		{"NaN lat", math.NaN(), 0, true},
		{"NaN lon", 0, math.NaN(), true},
		{"Inf lat", math.Inf(1), 0, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := NewCoordinateFromDeg(c.lat, c.lon)
			if c.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !c.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestNewCoordinateFromString(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		wantErr bool
	}{
		{"typical", `55°45'21" 37°37'2"`, false},
		{"multiple spaces", `55°45'21"    37°37'2"`, false},
		{"leading/trailing spaces", `  55°45'21" 37°37'2"  `, false},
		{"negative", `-55°45'21" -37°37'2"`, false},
		{"empty", "", true},
		{"one value", `55°45'21"`, true},
		{"three values", `55°45'21" 37°37'2" 10°0'0"`, true},
		{"bad lat", `91°0'0" 37°37'2"`, true},
		{"bad lon", `55°45'21" 181°0'0"`, true},
		{"incomplete lat", `55°45' 37°37'2"`, true},
		{"incomplete lon", `55°45'21" 37°37'`, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := NewCoordinateFromString(c.in)
			if c.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !c.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// Проверяем, что lat и lon не перепутаны.
func TestNewCoordinateFromStringValues(t *testing.T) {
	c, err := NewCoordinateFromString(`10°0'0" 20°0'0"`)
	if err != nil {
		t.Fatal(err)
	}
	lat, _ := c.latitude.GetDegree()
	lon, _ := c.longitude.GetDegree()
	if math.Abs(lat-10) > 1e-9 {
		t.Errorf("lat = %v, want 10", lat)
	}
	if math.Abs(lon-20) > 1e-9 {
		t.Errorf("lon = %v, want 20", lon)
	}
}

func TestCoordinateCheckValid(t *testing.T) {
	c, _ := NewCoordinateFromDeg(10, 20)
	if err := c.CheckValid(); err != nil {
		t.Errorf("valid coordinate returned error: %v", err)
	}
}

func TestCoordinateEqual(t *testing.T) {
	a, _ := NewCoordinateFromDeg(10, 20)
	b, _ := NewCoordinateFromDeg(10, 20)
	c, _ := NewCoordinateFromDeg(10, 20.001)

	if !a.Equal(b) {
		t.Error("identical coordinates should be equal")
	}
	if a.Equal(c) {
		t.Error("different coordinates should not be equal")
	}
}

func TestCoordinateEqualEpsilon(t *testing.T) {
	a, _ := NewCoordinateFromDeg(10, 20)
	b, _ := NewCoordinateFromDeg(10+1e-12, 20+1e-12)
	if !a.Equal(b) {
		t.Error("coordinates within eps should be equal")
	}
}
