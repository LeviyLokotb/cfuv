package steg

import "github.com/LeviyLokotb/imgsteg/internal/img"

type cell struct {
	// Индекс пиксела
	Px int
	// Индекс канала
	Ch img.ColorChan
	// Сколько бит доступно для записи
	Bits int
}

func buildCells(pxTotal int, bitsToEncode int, chansToEncode []img.ColorChan) []cell {
	cells := make([]cell, pxTotal*len(chansToEncode))

	ich := 0
	ipx := 0
	for i := range cells {
		cells[i] = cell{
			Px:   ipx,
			Ch:   chansToEncode[ich],
			Bits: bitsToEncode,
		}

		ich += 1
		if ich >= len(chansToEncode) {
			ich = 0
			ipx += 1
		}
	}
	return cells
}
