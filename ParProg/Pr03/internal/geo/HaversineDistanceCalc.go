package geo

import "fmt"

type HaversineDistanceCalc struct{}

func (dc HaversineDistanceCalc) DistanceBetween(a, b Coordinate) (float64, error) {
	d, err := haversine(a, b)
	if err != nil {
		return 0.0, fmt.Errorf("can not calculate distance: %w", err)
	}
	return d, nil
}
