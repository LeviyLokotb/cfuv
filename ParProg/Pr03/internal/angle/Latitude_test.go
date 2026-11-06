package angle

import (
	"math"
	"testing"
)

func TestLatitudeFromDeg(t *testing.T) {
	cases := []struct {
		name    string
		in      float64
		wantErr bool
	}{
		{"min", -90, false},
		{"max", 90, false},
		{"zero", 0, false},
		{"typical", 55.7558, false},
		{"below min", -90.1, true},
		{"above max", 90.1, true},
		{"NaN", math.NaN(), true},
		{"+Inf", math.Inf(1), true},
		{"-Inf", math.Inf(-1), true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := LatitudeFromDeg(c.in)
			if c.wantErr {
				if err == nil {
					t.Fatalf("FromDeg(%v) = %v, want error", c.in, got)
				}
				if got != nil {
					t.Errorf("on error expected nil angle, got %v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("FromDeg(%v) unexpected error: %v", c.in, err)
			}
			deg, err := got.GetDegree()
			if err != nil {
				t.Fatalf("GetDegree() error: %v", err)
			}
			if deg != c.in {
				t.Errorf("got %v, want %v", deg, c.in)
			}
		})
	}
}

func TestLatitudeFromString(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		want    float64
		wantErr bool
	}{
		{"typical", `55°45'21"`, 55 + 45.0/60 + 21.0/3600, false},
		{"negative", `-55°45'21"`, -(55 + 45.0/60 + 21.0/3600), false},
		{"zero", `0°0'0"`, 0, false},
		{"max", `90°0'0"`, 90, false},
		{"over max", `91°0'0"`, 0, true},
		{"empty", "", 0, true},
		{"garbage", "abc", 0, true},
		{"incomplete DMS", `10°30'`, 0, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := LatitudeFromString(c.in)
			if c.wantErr {
				if err == nil {
					t.Fatalf("FromString(%q) = %v, want error", c.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("FromString(%q) unexpected error: %v", c.in, err)
			}
			deg, _ := got.GetDegree()
			if math.Abs(deg-c.want) > 1e-9 {
				t.Errorf("got %v, want %v", deg, c.want)
			}
		})
	}
}

func TestLatitudeGetRadian(t *testing.T) {
	// Невалидная широта — FromDeg возвращает nil + error
	l, err := LatitudeFromDeg(180)
	if err == nil {
		t.Fatal("expected error for invalid latitude")
	}
	if l != nil {
		t.Errorf("expected nil angle, got %v", l)
	}

	// Валидная
	l2, err := LatitudeFromDeg(90)
	if err != nil {
		t.Fatal(err)
	}
	r, err := l2.GetRadian()
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(r-math.Pi/2) > 1e-12 {
		t.Errorf("got %v, want %v", r, math.Pi/2)
	}
}

func TestLatitudeEqual(t *testing.T) {
	a, _ := LatitudeFromDeg(10.0)
	b, _ := LatitudeFromDeg(10.0 + AngleEqualEps/2)
	c, _ := LatitudeFromDeg(10.0 + AngleEqualEps*10)

	if !a.Equal(b) {
		t.Error("values within eps should be equal")
	}
	if a.Equal(c) {
		t.Error("values outside eps should not be equal")
	}
}
