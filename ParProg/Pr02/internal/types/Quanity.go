package types

import "errors"

type Quanity int

const (
	QuanityInvalid Quanity = -1
	QuanityZero    Quanity = 0
)

func (a Quanity) Get() (int, error) {
	if a < 0 {
		return int(QuanityInvalid), errors.New("quanity must be > 0")
	} else if a == 0 {
		return int(QuanityZero), errors.New("quanity must be non-zero")
	}
	return int(a), nil
}

func (a Quanity) Add(b Quanity) Quanity {
	if a < 0 || b < 0 {
		return QuanityInvalid
	}
	return a + b
}
