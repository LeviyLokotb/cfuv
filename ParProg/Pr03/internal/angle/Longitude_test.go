package angle

import (
	"math"
	"testing"
)

func TestLongitudeFromDeg(t *testing.T) {
	cases := []struct {
		name    string
		in      float64
		wantErr bool
	}{
		{"min", -180, false},
		{"max", 180, false},
		{"zero", 0, false},
		{"typical", 37.6173, false},
		{"below min", -180.1, true},
		{"above max", 180.1, true},
		{"NaN", math.NaN(), true},
		{"+Inf", math.Inf(1), true},
		{"-Inf", math.Inf(-1), true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := LongitudeFromDeg(c.in)
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
			deg, _ := got.GetDegree()
			if deg != c.in {
				t.Errorf("got %v, want %v", deg, c.in)
			}
		})
	}
}

func TestLongitudeFromString(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		want    float64
		wantErr bool
	}{
		{"typical", `37°37'2"`, 37 + 37.0/60 + 2.0/3600, false},
		{"negative", `-37°37'2"`, -(37 + 37.0/60 + 2.0/3600), false},
		{"max", `180°0'0"`, 180, false},
		{"over max", `181°0'0"`, 0, true},
		{"empty", "", 0, true},
		{"incomplete DMS", `37°37'`, 0, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := LongitudeFromString(c.in)
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

func TestLongitudeEqual(t *testing.T) {
	a, _ := LongitudeFromDeg(10.0)
	b, _ := LongitudeFromDeg(10.0 + AngleEqualEps/2)
	c, _ := LongitudeFromDeg(10.0 + AngleEqualEps*10)

	if !a.Equal(b) {
		t.Error("values within eps should be equal")
	}
	if a.Equal(c) {
		t.Error("values outside eps should not be equal")
	}
}
