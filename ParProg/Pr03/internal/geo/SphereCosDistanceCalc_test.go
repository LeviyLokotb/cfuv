package geo

import (
	"math"
	"testing"
)

func TestSphereCosDistanceCalc(t *testing.T) {
	dc := SphereCosDistanceCalc{}
	a := mustCoord(t, 55.7558, 37.6173)
	b := mustCoord(t, 59.9343, 30.3351)

	d, err := dc.DistanceBetween(a, b)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(d-635_000) > 10_000 {
		t.Errorf("got %v, want ~635000", d)
	}
}

func TestSphereCosDistanceCalcSamePoint(t *testing.T) {
	dc := SphereCosDistanceCalc{}
	a := mustCoord(t, 10, 20)
	d, err := dc.DistanceBetween(a, a)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(d) > 1e-3 {
		t.Errorf("distance to self = %v, want ~0", d)
	}
}
