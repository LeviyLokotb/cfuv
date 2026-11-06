package steg

import (
	"github.com/LeviyLokotb/imgsteg/internal/img"
)

func lastByteIndex(bits int) int {
	return bits >> 3
}

// Читает bitsToRead бит из src с позиции bitOff, не больше 8 бит
//
// Возвращает байт с прочитанными битами в младших позициях
func readNBits(src []byte, bitOff int, bitsToRead int) byte {
	// Индекс байта
	ibyte := lastByteIndex(bitOff)
	// Отступ внутри нужного байта
	localOff := uint(bitOff & 0b111)

	if ibyte >= len(src) {
		return 0
	}

	// Склеиваем байты с информацией в 2-байтовое слово для удобства
	h := src[ibyte]
	l := byte(0)
	if ibyte+1 < len(src) {
		l = src[ibyte+1]
	}
	pair := uint16(h)<<8 + uint16(l)

	// Выравниваем нужную информацию
	pair <<= localOff

	// Двигаем значимые биты в младшие позиции результата ии обрезаем
	return byte(pair >> (16 - bitsToRead))
}

// Записывает bitsToWrite бит содержащихся в младших позициях bits
// в dst на позицию bitOff
func writeNBits(dst []byte, bits byte, bitOff int, bitsToWrite int) {
	// Индекс байта
	ibyte := lastByteIndex(bitOff)
	// Отступ внутри нужного байта
	localOff := uint(bitOff & 0b111)

	if ibyte >= len(dst) {
		return
	}

	// Склеиваем байты с информацией в 2-байтовое слово для удобства
	h := dst[ibyte]
	l := byte(0)
	if ibyte+1 < len(dst) {
		l = dst[ibyte+1]
	}
	pair := uint16(h)<<8 + uint16(l)

	// Строим изменённое слово
	// Слово с битами инфорамции, старший байт 0
	val := uint16(bits) << (16 - byte(bitsToWrite) - byte(localOff))
	// Маска для слова, показывает где должны быть биты информации
	valMask := uint16((^(byte(0xFF) >> bitsToWrite)) >> localOff)

	pair = (pair & ^valMask) | val

	h = byte(pair >> 8)
	l = byte(pair)

	dst[ibyte] = h
	if ibyte+1 < len(dst) {
		dst[ibyte+1] = l
	}
}

// Записывает bitsToEncode младших бит из val в dst байт через XOR
//
// Возвращает изменённый байт
func encodeNBits(dst byte, val byte, bitsToEncode int) byte {
	// Зашифрованные биты
	encByte := (dst ^ val) & ^(0xFF << bitsToEncode)

	// Освобождаем место в исходном байте
	clearDst := dst & (0xFF << bitsToEncode)

	return clearDst | encByte
}

// Декодирует bitsToDecode младших бит из src и orig байт через XOR
func decodeNBits(src byte, orig byte, bitsToDecode int) byte {
	// XOR с оригиналом даёт исходное значение младших бит
	return (src ^ orig) & ^(0xFF << bitsToDecode)
}

func MaxBitsToMode(bmp img.Bitmap, cfg Options) int {
	cells := buildCells(bmp.LenPx(), cfg.BitsToMod, cfg.ChansToEncode)
	maxBits := 0
	for _, cl := range cells {
		maxBits += cl.Bits
	}
	return maxBits
}

// Конвертирует число в байты в little-endian
func itob(n int) []byte {
	var bytes []byte
	for n != 0 {
		b := byte(n)
		bytes = append(bytes, b)
		n >>= 8
	}
	if len(bytes) == 0 {
		bytes = append(bytes, 0x00)
	}
	return bytes
}

// Конвертирует байты в little-endian в число
func btoi(bb []byte) int {
	n := 0
	for i, b := range bb {
		n += int(b) << (8 * i)
	}
	return n
}

// Добавляет число в виде 4-байтового префикса (little-endian) в начало
func withNumPrefix(src []byte, n uint32) []byte {
	bb := itob(int(n))
	// Дополняем до 4 байт
	for len(bb) < 4 {
		bb = append(bb, 0)
	}
	return append(bb, src...)
}
