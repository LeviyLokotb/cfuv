package geo

import (
	"errors"
	"fmt"
	"slices"
)

/*
Маршрут из последовательных точек
*/
type Route struct {
	positions          []Coordinate
	DistanceCalculator IDistanceCalc
	valid_state        validate_state
	last_err           error
}

func NewRoute(positions []Coordinate, dc IDistanceCalc) (Route, error) {
	r := Route{
		positions:          positions,
		DistanceCalculator: dc,
		valid_state:        state_unknown,
	}
	err := r.CheckValid()
	return r, err
}

func (r *Route) CheckValid() (_err error) {
	defer func() {
		r.last_err = _err
	}()

	switch r.valid_state {
	case state_valid:
		return nil
	case state_invalid:
		return r.last_err
	}

	if len(r.positions) < 2 {
		return errors.New("invalid route: length must be >= 2")
	}
	for i := 0; i < len(r.positions)-1; i++ {
		a, b := r.positions[i], r.positions[i+1]

		if err := a.CheckValid(); err != nil {
			r.valid_state = state_invalid
			return fmt.Errorf("invalid route: %w", err)
		}
		if err := b.CheckValid(); err != nil {
			r.valid_state = state_invalid
			return fmt.Errorf("invalid route: %w", err)
		}

		if a.Equal(b) {
			r.valid_state = state_invalid
			return errors.New("invalid route: two identical values in a row")
		}
	}
	r.valid_state = state_valid
	return nil
}

func (r *Route) AddPoint(c Coordinate) error {
	nr := *r

	// Используем новый массив, потому что там сейчас указатель
	nr.positions = make([]Coordinate, len(r.positions))
	copy(nr.positions, r.positions)

	nr.positions = append(nr.positions, c)
	nr.valid_state = state_unknown

	err := nr.CheckValid()
	if err != nil {
		return err
	}
	*r = nr
	return nil
}

// Длина маршрута в метрах
//
// Если маршрут невалидный, возвращает 0
func (r Route) Length() (float64, error) {
	if err := r.CheckValid(); err != nil {
		return 0.0, err
	}
	length := 0.0
	for i := 0; i < len(r.positions)-1; i++ {
		a, b := r.positions[i], r.positions[i+1]

		d, err := r.DistanceCalculator.DistanceBetween(a, b)
		if err != nil {
			return 0.0, err
		}
		length += d
	}
	return length, nil
}

// Расстояние между точкой и следующей за ней в метрах
func (r Route) GetDistanceAfterIdx(idx int) (float64, error) {
	if (idx < 0) || (idx+1 >= len(r.positions)) {
		return 0.0, errors.New("invalid point index")
	}
	d, err := r.DistanceCalculator.DistanceBetween(r.positions[idx], r.positions[idx+1])
	if err != nil {
		return 0.0, fmt.Errorf("cannot calculate distance: %w", err)
	}
	return d, nil
}

// Расстояние между точкой и следующей за ней в метрах
func (r Route) GetDistanceAfter(c Coordinate) (float64, error) {
	idx := slices.IndexFunc(r.positions, func(p Coordinate) bool {
		return c.Equal(p)
	})
	if idx == -1 {
		return 0.0, errors.New("invalid coordinate")
	}
	return r.GetDistanceAfterIdx(idx)
}
