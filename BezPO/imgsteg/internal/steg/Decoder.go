package steg

import (
	"errors"
	"fmt"

	"github.com/LeviyLokotb/imgsteg/internal/img"
)

type Decoder struct {
	imageBase img.Bitmap
	o         Options
	// Виртуальные ячейки памяти (содержат информацию о доступных для записи битах)
	cells []cell
}

func NewDecoder(baseImg img.Bitmap, o Options) (Decoder, error) {
	if err := o.Validate(); err != nil {
		return Decoder{}, fmt.Errorf("failed to create encoder: %w", err)
	}
	if baseImg.Px == nil {
		return Decoder{}, errors.New("failed to create encoder: bad image")
	}
	return Decoder{
		imageBase: baseImg,
		o:         o,
		cells:     buildCells(baseImg.LenPx(), o.BitsToMod, o.ChansToEncode),
	}, nil
}

func (d Decoder) DecodeBytes(imageEnc img.Bitmap) (MetaData, error) {
	if imageEnc.Px == nil {
		return MetaData{}, errors.New("encoded image is empty")
	}
	if imageEnc.W != d.imageBase.W || imageEnc.H != d.imageBase.H {
		return MetaData{}, fmt.Errorf(
			"image size mismatch: base %dx%d, encoded %dx%d",
			d.imageBase.W, d.imageBase.H, imageEnc.W, imageEnc.H,
		)
	}
	if len(imageEnc.Px) < len(d.imageBase.Px) {
		return MetaData{}, errors.New("encoded image is smaller than base")
	}

	maxBits := MaxBitsToMode(d.imageBase, d.o)
	if maxBits < 32 { // 4 байта префикса
		return MetaData{}, errors.New("image too small to contain a message")
	}
	maxBytes := lastByteIndex(maxBits) + 1
	if maxBytes < 4 {
		return MetaData{}, errors.New("image too small to contain a message")
	}

	result := make([]byte, maxBytes)
	ibit := 0
	for _, cl := range d.cells {
		// Читаем пока есть что читать
		if ibit >= maxBits {
			break
		}

		src := imageEnc.Px[cl.Px][cl.Ch]
		org := d.imageBase.Px[cl.Px][cl.Ch]
		v := decodeNBits(src, org, cl.Bits)
		// Записываем в результат
		writeNBits(result, v, ibit, cl.Bits)

		ibit += cl.Bits
	}

	data, err := d.o.Prep.DecodeMeta(result)
	if err != nil {
		return MetaData{}, err
	}

	return data, nil
}
