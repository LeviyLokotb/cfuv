package geo

import "fmt"

type SphereCosDistanceCalc struct{}

func (dc SphereCosDistanceCalc) DistanceBetween(a, b Coordinate) (float64, error) {
	d, err := sphereCosTheorem(a, b)
	if err != nil {
		return 0.0, fmt.Errorf("can not calculate distance: %w", err)
	}
	return d, nil
}
