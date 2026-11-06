package angle

import (
	"errors"
	"math"
)

/*
Долгота. Реализует IAngle, включает значения [-180;180].
*/
type Longitude float64

const (
	LongitudeMinDeg   Longitude = -180.0
	LongitudeMaxDeg   Longitude = 180.0
	LongitudeInstance Longitude = 0.0
)

// Создание объекта Longitude из числа.
// Гарантируется соблюдение диапазона
func LongitudeFromDeg(deg float64) (IAngle, error) {
	l := Longitude(deg)
	if _, err := l.GetDegree(); err != nil {
		return nil, err
	}
	return l, nil
}

/*
Создание объекта Longitude из строки
формата Градусы-Минуты-Секунды.

Строка имеет формат:

D.d°M.m'S.s"

Обозначением градусов служат: ° (U+00B0)

Обозначением минут служат: ' (U+0027) ` (U+0060) ′ (U+2032)

Обозначением секунд служат: " (U+0022) ″ (U+2033) ' (U+0027)x2 ′ (U+2032)x2 ` (U+0060)x2

Гарантируется соблюдение диапазона
*/
func LongitudeFromString(s string) (IAngle, error) {
	d, err := parseDMS(s)
	if err != nil {
		return nil, err
	}
	return LongitudeFromDeg(d)
}

func (l Longitude) GetDegree() (float64, error) {
	if math.IsNaN(float64(l)) || l < LongitudeMinDeg || l > LongitudeMaxDeg {
		return 0, errors.New("invalid longitude")
	}
	return float64(l), nil
}

func (l Longitude) GetRadian() (float64, error) {
	a, err := l.GetDegree()
	return DegToRad(a), err
}

// Равны ли углы. Возвращает false если угол некорректный
func (l Longitude) Equal(a IAngle) bool {
	nl, err := l.GetDegree()
	if err != nil {
		return false
	}
	na, err := a.GetDegree()
	if err != nil {
		return false
	}
	return floatEqual(nl, na)
}
