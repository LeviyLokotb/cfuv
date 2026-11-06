package angle

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

func DegToRad(d float64) float64 {
	return d / 180.0 * math.Pi
}

func RadToDeg(r float64) float64 {
	return r * 180.0 / math.Pi
}

func dmsToDeg(d, m, s float64) float64 {
	return d + m/60.0 + s/3600.0
}

const (
	signRe   = `([+-])?`
	spaceRe  = `\s*`
	numRe    = `(\d+(?:[.]\d+)?)`           // Число с возможной мантисой
	degDelRe = `(?:°)`                      // Разделитель градусов (°)
	minDelRe = `(?:'|′|` + "`" + `)`        // Разделитель минут    (' ′ `)
	secDelRe = `(?:''|"|′′|″|` + "``" + `)` // Разделитель секунд   ('' " ′′ ″ `)
)

var dmsReg = regexp.MustCompile(
	`^` + spaceRe +
		signRe + numRe + degDelRe + spaceRe +
		numRe + minDelRe + spaceRe +
		numRe + secDelRe +
		spaceRe + `$`,
)

func parse(s string) (float64, error) {
	if s == "" {
		return 0, nil
	}
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid DMS format: %w", err)
	}
	return n, nil
}

func parseDMS(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, errors.New("empty string")
	}
	s = strings.ReplaceAll(s, ",", ".")

	m := dmsReg.FindStringSubmatch(s)
	if m == nil {
		return 0, fmt.Errorf("invalid DMS format: %q", s)
	}

	deg, err := parse(m[2])
	if err != nil {
		return 0, err
	}
	min, err := parse(m[3])
	if err != nil {
		return 0, err
	}
	sec, err := parse(m[4])
	if err != nil {
		return 0, err
	}
	sign := 1.0
	if m[1] == "-" {
		sign = -1.0
	}

	return sign * dmsToDeg(deg, min, sec), nil
}

const AngleEqualEps = 1e-9 // точность порядка 3 микрон на местности

func floatEqual(a, b float64) bool {
	return math.Abs(a-b) < AngleEqualEps
}
