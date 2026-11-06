package steg

import (
	"errors"
	"fmt"
	"io"

	"github.com/LeviyLokotb/imgsteg/internal/img"
)

type Encoder struct {
	imageBase img.Bitmap
	imageEnc  img.Bitmap
	o         Options
	// Виртуальные ячейки памяти (содержат информацию о доступных для записи битах)
	cells []cell
}

func NewEncoder(baseImg img.Bitmap, o Options) (Encoder, error) {
	if err := o.Validate(); err != nil {
		return Encoder{}, fmt.Errorf("failed to create encoder: %w", err)
	}
	if baseImg.Px == nil {
		return Encoder{}, errors.New("failed to create encoder: bad image")
	}
	return Encoder{
		imageBase: baseImg,
		imageEnc:  baseImg.Clone(),
		o:         o,
		cells:     buildCells(baseImg.LenPx(), o.BitsToMod, o.ChansToEncode),
	}, nil
}

// Записывает байты data в изображение.
// Повторный вызов перезаписывает информацию а не добавляет в конец
//
// Возвращает io.EOF при попытке записи за границы
func (e *Encoder) EncodeBytes(d MetaData) error {
	// Добавляем префикс с длиной сообщения
	data := e.o.Prep.EncodeMeta(d)

	// доступно бит в ячейках для записи
	maxBits := MaxBitsToMode(e.imageBase, e.o)

	// всего нужно записать бит
	bitsTotal := len(data) * 8
	if bitsTotal > maxBits {
		return errors.New("oversized data")
	}

	// Обнуляем результат
	e.imageEnc = e.imageBase.Clone()

	// счётчик записанных битов
	ibit := 0
	for _, cl := range e.cells {
		// Пишем пока есть что писать
		if ibit >= bitsTotal {
			break
		}
		// Если вышли за пределы
		if ibit >= maxBits {
			return io.EOF
		}
		// Читаем столько бит, сколько поместится в текущий слот
		nextBits := readNBits(data, ibit, cl.Bits)

		orgByte := e.imageBase.Px[cl.Px][cl.Ch]
		encByte := encodeNBits(orgByte, nextBits, cl.Bits)
		e.imageEnc.Px[cl.Px][cl.Ch] = encByte

		ibit += cl.Bits
	}

	return nil
}

func (e Encoder) GetEncoded() img.Bitmap {
	return e.imageEnc.Clone()
}
