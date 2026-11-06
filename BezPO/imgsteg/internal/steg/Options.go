package steg

import (
	"errors"

	"github.com/LeviyLokotb/imgsteg/internal/img"
)

type Options struct {
	// Сколько можно изменить битов в канале
	BitsToMod int
	// Какие каналы участвуют в кодировке
	ChansToEncode []img.ColorChan

	Prep Preprocessor
}

func NewOptions(bitsToEncode int, chansToEncode []img.ColorChan) (Options, error) {
	o := Options{
		BitsToMod:     bitsToEncode,
		ChansToEncode: chansToEncode,
	}
	return o, o.Validate()
}

func NewOptionsMinimal() Options {
	o, _ := NewOptions(1, []img.ColorChan{img.Red})
	return o
}

func NewOptionsOptimal() Options {
	o, _ := NewOptions(1, []img.ColorChan{img.Red, img.Green, img.Blue})
	return o
}

func NewOptionsRedChan() Options {
	o, _ := NewOptions(3, []img.ColorChan{img.Red})
	return o
}

func NewOptionsAgressive() Options {
	o, _ := NewOptions(3, []img.ColorChan{img.Red, img.Green, img.Blue})
	return o
}

func (o Options) Validate() error {
	if len(o.ChansToEncode) > 3 || len(o.ChansToEncode) < 1 {
		return errors.New("invalid config: 1 <= chans <= 3")
	}
	if o.BitsToMod > 8 || o.BitsToMod < 1 {
		return errors.New("invalid config: 1 <= bits <= 8")
	}
	return nil
}
