package angle

import (
	"math"
	"testing"
)

func TestDegToRad(t *testing.T) {
	cases := []struct {
		deg  float64
		want float64
	}{
		{0, 0},
		{180, math.Pi},
		{-90, -math.Pi / 2},
		{360, 2 * math.Pi},
	}
	for _, c := range cases {
		got := DegToRad(c.deg)
		if math.Abs(got-c.want) > 1e-12 {
			t.Errorf("DegToRad(%v) = %v, want %v", c.deg, got, c.want)
		}
	}
}

func TestRadToDeg(t *testing.T) {
	cases := []struct {
		rad  float64
		want float64
	}{
		{0, 0},
		{math.Pi, 180},
		{-math.Pi / 2, -90},
	}
	for _, c := range cases {
		got := RadToDeg(c.rad)
		if math.Abs(got-c.want) > 1e-9 {
			t.Errorf("RadToDeg(%v) = %v, want %v", c.rad, got, c.want)
		}
	}
}

func TestDegRadRoundTrip(t *testing.T) {
	for _, deg := range []float64{-180, -90, -45, 0, 45, 90, 180} {
		got := RadToDeg(DegToRad(deg))
		if math.Abs(got-deg) > 1e-9 {
			t.Errorf("round trip %v -> %v", deg, got)
		}
	}
}

// dmsToDeg — чистая арифметика: знак применяется только к градусам.
// Знак всей DMS-записи обрабатывается снаружи в parseDMS.
func TestDmsToDeg(t *testing.T) {
	cases := []struct {
		d, m, s float64
		want    float64
	}{
		{0, 0, 0, 0},
		{10, 30, 0, 10.5},
		{10, 0, 30, 10 + 30.0/3600},
		{55, 45, 21, 55 + 45.0/60 + 21.0/3600},
		{-10, 30, 0, -9.5},
		{-10, -30, 0, -10.5}, // если оба отрицательны — складываются как есть
	}
	for _, c := range cases {
		got := dmsToDeg(c.d, c.m, c.s)
		if math.Abs(got-c.want) > 1e-12 {
			t.Errorf("dmsToDeg(%v,%v,%v) = %v, want %v", c.d, c.m, c.s, got, c.want)
		}
	}
}

// Текущий контракт: строгий формат D°M'S" — все три компонента обязательны.
// Усечённые формы (10°, 10°30') отвергаются.
func TestParseDMS(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		want    float64
		wantErr bool
	}{
		{"deg-min-sec double quote", `10°30'15"`, 10 + 30.0/60 + 15.0/3600, false},
		{"deg-min-sec unicode prime", "10°30′15″", 10 + 30.0/60 + 15.0/3600, false},
		{"backtick minute and second", "10°30`15``", 10 + 30.0/60 + 15.0/3600, false},
		{"negative", `-10°30'15"`, -(10 + 30.0/60 + 15.0/3600), false},
		{"positive sign", `+10°30'15"`, 10 + 30.0/60 + 15.0/3600, false},
		{"comma as decimal separator", `10°30'15,5"`, 10 + 30.0/60 + 15.5/3600, false},
		{"fractional seconds", `10°30'15.5"`, 10 + 30.0/60 + 15.5/3600, false},
		{"spaces around", `  10° 30' 15"  `, 10 + 30.0/60 + 15.0/3600, false},
		{"fractional degrees", `10.5°30'15"`, 10.5 + 30.0/60 + 15.0/3600, false},

		// Строгий формат: усечённые записи не принимаются
		{"plain degrees rejected", "10°", 0, true},
		{"deg-min rejected", "10°30'", 0, true},

		{"empty", "", 0, true},
		{"only spaces", "   ", 0, true},
		{"garbage", "abc", 0, true},
		{"minutes without degrees", `30'15"`, 0, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := parseDMS(c.in)
			if c.wantErr {
				if err == nil {
					t.Fatalf("parseDMS(%q) = %v, want error", c.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseDMS(%q) unexpected error: %v", c.in, err)
			}
			if math.Abs(got-c.want) > 1e-9 {
				t.Errorf("parseDMS(%q) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}

func TestFloatEqual(t *testing.T) {
	if !floatEqual(1.0, 1.0) {
		t.Error("equal values should be equal")
	}
	if !floatEqual(1.0, 1.0+AngleEqualEps/2) {
		t.Error("values within eps should be equal")
	}
	if floatEqual(1.0, 1.0+AngleEqualEps*2) {
		t.Error("values outside eps should not be equal")
	}
}

// FuzzParseDMS проверяет, что парсер не паникует на произвольном вводе.
func FuzzParseDMS(f *testing.F) {
	f.Add(`10°30'15"`)
	f.Add(`-10°30'15"`)
	f.Add(`10,5°30'15.5"`)
	f.Add(`10°30′15″`)
	f.Add(``)
	f.Add(`abc`)
	f.Fuzz(func(t *testing.T, s string) {
		_, _ = parseDMS(s)
	})
}
