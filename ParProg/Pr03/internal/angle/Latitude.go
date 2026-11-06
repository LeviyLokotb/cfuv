package angle

import (
	"errors"
	"math"
)

/*
Широта. Реализует IAngle, включает значения [-90;90].
*/
type Latitude float64

const (
	LatitudeMinDeg   Latitude = -90.0
	LatitudeMaxDeg   Latitude = 90.0
	LatitudeInstance Latitude = 0.0
)

// Создание объекта Latitude из числа.
// Гарантируется соблюдение диапазона
func LatitudeFromDeg(deg float64) (IAngle, error) {
	l := Latitude(deg)
	if _, err := l.GetDegree(); err != nil {
		return nil, err
	}
	return l, nil
}

/*
Создание объекта Latitude из строки
формата Градусы-Минуты-Секунды.

Строка имеет формат:

D.d°M.m'S.s"

Обозначением градусов служат следующие символы: ° (U+00B0)

Обозначением минут служат следующие символы: ' (U+0027) ` (U+0060) ′ (U+2032)

Обозначением секунд служат следующие символы: " (U+0022) ″ (U+2033) ' (U+0027)x2 ′ (U+2032)x2 `(U+0060)x2

Гарантируется соблюдение диапазона
*/
func LatitudeFromString(s string) (IAngle, error) {
	d, err := parseDMS(s)
	if err != nil {
		return nil, err
	}
	return LatitudeFromDeg(d)
}

func (l Latitude) GetDegree() (float64, error) {
	if math.IsNaN(float64(l)) || l < LatitudeMinDeg || l > LatitudeMaxDeg {
		return 0, errors.New("invalid latitude")
	}
	return float64(l), nil
}

func (l Latitude) GetRadian() (float64, error) {
	a, err := l.GetDegree()
	return DegToRad(a), err
}

// Равны ли углы. Возвращает false если угол некорректный
func (l Latitude) Equal(a IAngle) bool {
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
