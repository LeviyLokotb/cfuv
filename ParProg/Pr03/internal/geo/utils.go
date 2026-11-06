package geo

import "math"

type validate_state uint8

const (
	state_unknown validate_state = 0
	state_valid   validate_state = 1
	state_invalid validate_state = 2
)

const (
	earth_radius_m = 6_371_000
)

// Расстояние между точками заданными координатами
//
// (формула гаверсинусов)
func haversine(a Coordinate, b Coordinate) (float64, error) {
	lat1, err := a.latitude.GetRadian()
	if err != nil {
		return 0, err
	}

	lon1, err := a.longitude.GetRadian()
	if err != nil {
		return 0, err
	}

	lat2, err := b.latitude.GetRadian()
	if err != nil {
		return 0, err
	}

	lon2, err := b.longitude.GetRadian()
	if err != nil {
		return 0, err
	}

	dlat := lat2 - lat1
	dlon := lon2 - lon1

	lat_sin := math.Sin(dlat / 2)
	lon_sin := math.Sin(dlon / 2)
	na := lat_sin*lat_sin + math.Cos(lat1)*math.Cos(lat2)*lon_sin*lon_sin
	if na > 1 {
		na = 1
	}
	if na < 0 {
		na = 0
	}

	nc := 2 * math.Atan2(math.Sqrt(na), math.Sqrt(1-na))

	return earth_radius_m * nc, nil
}

// Расстояние между точками заданными координатами
//
// (формула сферической теоремы косинусов)
func sphereCosTheorem(a Coordinate, b Coordinate) (float64, error) {
	lat1, err := a.latitude.GetRadian()
	if err != nil {
		return 0, err
	}

	lon1, err := a.longitude.GetRadian()
	if err != nil {
		return 0, err
	}

	lat2, err := b.latitude.GetRadian()
	if err != nil {
		return 0, err
	}

	lon2, err := b.longitude.GetRadian()
	if err != nil {
		return 0, err
	}

	dlon := lon2 - lon1

	s := math.Sin(lat1) * math.Sin(lat2)
	c := math.Cos(lat1) * math.Cos(lat2) * math.Cos(dlon)

	arg := s + c
	if arg > 1 {
		arg = 1
	}
	if arg < -1 {
		arg = -1
	}
	return earth_radius_m * math.Acos(arg), nil
}
