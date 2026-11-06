package geo

import (
	"errors"
	"fmt"
	"pr03/internal/angle"
	"strings"
)

/*
Географическая координата.

Включает Широту (angle.Latitude) и Долготу (angle.Longitude)
*/
type Coordinate struct {
	latitude  angle.IAngle
	longitude angle.IAngle
}

/*
Создание объекта Coordinate
из строки в формате Градусы-Минуты-Секунды

Ожидает две записи Градусы-Минуты-Секунды
(в формате соответствующем IAngle.FromString)
идущие через пробел

Гарантируется соблюдение диапазонов
*/
func NewCoordinateFromString(s string) (Coordinate, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Coordinate{}, errors.New("empty string")
	}

	l := strings.Fields(s)
	if len(l) != 2 {
		return Coordinate{}, errors.New("invalid coordinate format")
	}

	lat, err := angle.LatitudeFromString(l[0])
	if err != nil {
		return Coordinate{}, fmt.Errorf("invalid latitude format: %w", err)
	}
	lon, err := angle.LongitudeFromString(l[1])
	if err != nil {
		return Coordinate{}, fmt.Errorf("invalid longitude format: %w", err)
	}

	return Coordinate{
		latitude:  lat,
		longitude: lon,
	}, nil
}

/*
Создание объекта Coordinate из чисел (градусов)

Гарантируется соблюдение диапазонов
*/
func NewCoordinateFromDeg(lat, lon float64) (Coordinate, error) {
	clat, err := angle.LatitudeFromDeg(lat)
	if err != nil {
		return Coordinate{}, fmt.Errorf("invalid coordinate: %w", err)
	}

	clon, err := angle.LongitudeFromDeg(lon)
	if err != nil {
		return Coordinate{}, fmt.Errorf("invalid coordinate: %w", err)
	}

	return Coordinate{
		latitude:  clat,
		longitude: clon,
	}, nil
}

/*
Создание объекта Coordinate из углов пакета angle.

Соблюдение условий на ответственности программиста
*/
func NewCoordinateFromAngles(lat angle.IAngle, lon angle.IAngle) Coordinate {
	return Coordinate{
		latitude:  lat,
		longitude: lon,
	}
}

func (c Coordinate) CheckValid() error {
	if c.latitude == nil || c.longitude == nil {
		return errors.New("invalid coordinate: nil field")
	}

	if _, err := c.latitude.GetDegree(); err != nil {
		return fmt.Errorf("invalid coordinate: %w", err)
	}
	if _, err := c.longitude.GetDegree(); err != nil {
		return fmt.Errorf("invalid coordinate: %w", err)
	}
	return nil
}

// Проверка совпадения координат (точность angle.AngleEqualEps)
func (c Coordinate) Equal(p Coordinate) bool {
	return c.latitude.Equal(p.latitude) && c.longitude.Equal(p.longitude)
}
