package steg

// import (
// 	"bytes"
// 	"errors"
// 	"io"
// 	"math/rand"
// 	"testing"

// 	"github.com/LeviyLokotb/imgsteg/internal/img"
// )

// func Test_EncodeDecode_Optimal(t *testing.T) {
// 	opt := NewOptionsOptimal()

// 	// 3*100 = 300 бит для записи (37 байт)
// 	baseImg := img.NewBitmap(10, 10)

// 	enc, _ := NewEncoder(baseImg, opt)
// 	dec, _ := NewDecoder(baseImg, opt)

// 	sample := []byte("AB")

// 	err := enc.EncodeBytes(sample)
// 	if err != nil {
// 		t.Fatalf("WriteBytes error: %s", err)
// 	}
// 	imageEnc := enc.GetEncoded()

// 	decoded, _ := dec.DecodeBytes(imageEnc)

// 	if len(sample) != len(decoded) {
// 		t.Fatalf("len not match (%d and %d)", len(sample), len(decoded))
// 	}
// 	for i := range sample {
// 		if sample[i] != decoded[i] {
// 			t.Errorf("bytes not match (%d and %d)", sample[i], decoded[i])
// 		}
// 	}
// }

// // config описывает один набор опций для теста.
// type config struct {
// 	name string
// 	opt  Options
// 	// сколько бит тратится на один пиксель (для расчёта ёмкости)
// 	bitsPerPixel int
// }

// func configs() []config {
// 	return []config{
// 		{"Minimal", NewOptionsMinimal(), 1},     // 1 бит × 1 канал
// 		{"Optimal", NewOptionsOptimal(), 3},     // 1 бит × 3 канала
// 		{"RedChan", NewOptionsRedChan(), 3},     // 3 бита × 1 канал
// 		{"Agressive", NewOptionsAgressive(), 9}, // 3 бита × 3 канала
// 	}
// }

// func Test_EncodeBytes_TooLarge(t *testing.T) {
// 	for _, cfg := range configs() {
// 		cfg := cfg
// 		t.Run(cfg.name, func(t *testing.T) {
// 			// ровно 10 пикселей ёмкости
// 			baseImg := img.NewBitmap(2, 5)
// 			capacityBits := 10 * cfg.bitsPerPixel
// 			// хотим записать больше, чем влезает (с учётом 32-битного префикса)
// 			tooBig := make([]byte, capacityBits/8+8)

// 			enc, err := NewEncoder(baseImg, cfg.opt)
// 			if err != nil {
// 				t.Fatalf("NewEncoder: %v", err)
// 			}

// 			err = enc.EncodeBytes(tooBig)
// 			if err == nil {
// 				t.Fatal("expected error for oversized data, got nil")
// 			}
// 		})
// 	}
// }

// func Test_EncodeBytes_ExactFit(t *testing.T) {
// 	for _, cfg := range configs() {
// 		cfg := cfg
// 		t.Run(cfg.name, func(t *testing.T) {
// 			// подбираем число пикселей так, чтобы ёмкость = 32 + N*8 ровно
// 			// (префикс 32 бита + N байт данных)
// 			// при bitsPerPixel = 1, 3, 9 — подбираем N и px вручную
// 			var px, n int
// 			switch cfg.bitsPerPixel {
// 			case 1:
// 				px, n = 40, 1 // 40 бит = 32 + 8
// 			case 3:
// 				// 3*px = 32 + 8*n → px = (32+8n)/3
// 				// n=2 → px=16 (48 бит = 32+16)
// 				px, n = 16, 2
// 			case 9:
// 				// 9*px = 32 + 8*n
// 				// n=2 → px=48/9 — не целое
// 				// n=5 → (32+40)=72, px=8 → ровно
// 				px, n = 8, 5
// 			default:
// 				t.Skipf("no exact-fit case for bitsPerPixel=%d", cfg.bitsPerPixel)
// 			}

// 			baseImg := img.NewBitmap(px, 1)
// 			sample := bytes.Repeat([]byte{0xA5}, n)

// 			enc, err := NewEncoder(baseImg, cfg.opt)
// 			if err != nil {
// 				t.Fatalf("NewEncoder: %v", err)
// 			}
// 			dec, err := NewDecoder(baseImg, cfg.opt)
// 			if err != nil {
// 				t.Fatalf("NewDecoder: %v", err)
// 			}

// 			if err := enc.EncodeBytes(sample); err != nil {
// 				t.Fatalf("EncodeBytes: %v", err)
// 			}

// 			got, _ := dec.DecodeBytes(enc.GetEncoded())
// 			if !bytes.Equal(got, sample) {
// 				t.Errorf("exact-fit mismatch: got %v, want %v", got, sample)
// 			}
// 		})
// 	}
// }

// // =====================================================================
// // utils.go
// // =====================================================================

// func Test_lastByteIndex(t *testing.T) {
// 	cases := []struct {
// 		bits int
// 		want int
// 	}{
// 		{0, 0}, {1, 0}, {7, 0}, {8, 1}, {9, 1}, {15, 1}, {16, 2}, {31, 3}, {32, 4},
// 	}
// 	for _, c := range cases {
// 		if got := lastByteIndex(c.bits); got != c.want {
// 			t.Errorf("lastByteIndex(%d) = %d, want %d", c.bits, got, c.want)
// 		}
// 	}
// }

// func Test_itob_btoi_RoundTrip(t *testing.T) {
// 	cases := []int{0, 1, 0xFF, 0x100, 0x123456, 0x7FFFFFFF}
// 	for _, n := range cases {
// 		b := itob(n)
// 		if got := btoi(b); got != n {
// 			t.Errorf("itob/btoi roundtrip %d -> %v -> %d", n, b, got)
// 		}
// 	}
// }

// func Test_itob_Zero(t *testing.T) {
// 	if got := itob(0); !bytes.Equal(got, []byte{0x00}) {
// 		t.Errorf("itob(0) = %v, want [0]", got)
// 	}
// }

// func Test_btoi_Empty(t *testing.T) {
// 	if got := btoi(nil); got != 0 {
// 		t.Errorf("btoi(nil) = %d, want 0", got)
// 	}
// }

// func Test_withPrefix(t *testing.T) {
// 	// 4 байта little-endian префикса + данные
// 	data := []byte("hello")
// 	got := withNumPrefix(data, uint32(len(data)))

// 	if len(got) != 4+len(data) {
// 		t.Fatalf("len(withPrefix) = %d, want %d", len(got), 4+len(data))
// 	}
// 	if btoi(got[:4]) != len(data) {
// 		t.Errorf("prefix = %d, want %d", btoi(got[:4]), len(data))
// 	}
// 	if !bytes.Equal(got[4:], data) {
// 		t.Errorf("data mismatch: got %v want %v", got[4:], data)
// 	}
// }

// func Test_withPrefix_Large(t *testing.T) {
// 	data := []byte("x")
// 	n := uint32(0x01020304)
// 	got := withNumPrefix(data, n)
// 	if btoi(got[:4]) != int(n) {
// 		t.Errorf("prefix = %d, want %d", btoi(got[:4]), n)
// 	}
// }

// // readNBits / writeNBits — round-trip для всех смещений и длин
// func Test_writeNBits_readNBits_RoundTrip(t *testing.T) {
// 	dst := make([]byte, 8)
// 	val := byte(0xAB) // 1010 1011

// 	for off := 0; off < 8*8-8; off++ {
// 		for n := 1; n <= 8; n++ {
// 			// очистим
// 			for i := range dst {
// 				dst[i] = 0
// 			}
// 			writeNBits(dst, val, off, n)
// 			got := readNBits(dst, off, n)
// 			want := val & byte(0xFF>>(8-n))
// 			if got != want {
// 				t.Errorf("off=%d n=%d: readNBits=%08b, want %08b", off, n, got, want)
// 			}
// 		}
// 	}
// }

// func Test_readNBits_OutOfRange(t *testing.T) {
// 	dst := []byte{0xFF, 0xFF}
// 	// читаем с позиции, где байт уже за пределами
// 	got := readNBits(dst, 8*len(dst), 4)
// 	if got != 0 {
// 		t.Errorf("readNBits beyond end = %d, want 0", got)
// 	}
// }

// func Test_writeNBits_OutOfRange(t *testing.T) {
// 	dst := []byte{0x00}
// 	// не должно паниковать
// 	writeNBits(dst, 0xFF, 8, 4)
// 	if dst[0] != 0x00 {
// 		t.Errorf("writeNBits beyond end modified buffer: %v", dst)
// 	}
// }

// func Test_encodeNBits_decodeNBits_RoundTrip(t *testing.T) {
// 	// для всех значений val (8 бит) и для всех длин 1..8
// 	for val := 0; val < 256; val++ {
// 		for n := 1; n <= 8; n++ {
// 			orig := byte(0xA5)
// 			enc := encodeNBits(orig, byte(val), n)
// 			dec := decodeNBits(enc, orig, n)
// 			want := byte(val) & byte(0xFF>>(8-n))
// 			if dec != want {
// 				t.Errorf("val=%08b n=%d: enc=%08b dec=%08b want %08b",
// 					val, n, enc, dec, want)
// 			}
// 		}
// 	}
// }

// func Test_encodeNBits_PreservesHighBits(t *testing.T) {
// 	orig := byte(0b11110000)
// 	enc := encodeNBits(orig, 0b00001111, 4)
// 	// старшие 4 бита должны остаться как в orig
// 	if enc&0b11110000 != 0b11110000 {
// 		t.Errorf("high bits not preserved: %08b", enc)
// 	}
// 	// младшие 4 бита — это val
// 	if enc&0b00001111 != 0b00001111 {
// 		t.Errorf("low bits wrong: %08b", enc)
// 	}
// }

// func Test_MaxBitsToMode(t *testing.T) {
// 	bmp := img.NewBitmap(4, 4) // 16 пикселей

// 	cases := []struct {
// 		name string
// 		opt  Options
// 		want int
// 	}{
// 		{"Minimal", NewOptionsMinimal(), 16 * 1 * 1},
// 		{"Optimal", NewOptionsOptimal(), 16 * 1 * 3},
// 		{"RedChan", NewOptionsRedChan(), 16 * 3 * 1},
// 		{"Agressive", NewOptionsAgressive(), 16 * 3 * 3},
// 	}
// 	for _, c := range cases {
// 		if got := MaxBitsToMode(bmp, c.opt); got != c.want {
// 			t.Errorf("%s: MaxBitsToMode = %d, want %d", c.name, got, c.want)
// 		}
// 	}
// }

// // =====================================================================
// // cell.go
// // =====================================================================

// func Test_buildCells_Count(t *testing.T) {
// 	chans := []img.ColorChan{img.Red, img.Green, img.Blue}
// 	cells := buildCells(4, 1, chans)
// 	if len(cells) != 4*3 {
// 		t.Fatalf("len(cells) = %d, want %d", len(cells), 4*3)
// 	}
// }

// func Test_buildCells_Order(t *testing.T) {
// 	chans := []img.ColorChan{img.Red, img.Green, img.Blue}
// 	cells := buildCells(3, 2, chans)

// 	// Ожидаемый порядок: (0,R),(0,G),(0,B),(1,R),(1,G),(1,B),(2,R),(2,G),(2,B)
// 	type expect struct {
// 		px int
// 		ch img.ColorChan
// 	}
// 	want := []expect{
// 		{0, img.Red}, {0, img.Green}, {0, img.Blue},
// 		{1, img.Red}, {1, img.Green}, {1, img.Blue},
// 		{2, img.Red}, {2, img.Green}, {2, img.Blue},
// 	}
// 	for i, w := range want {
// 		if cells[i].Px != w.px || cells[i].Ch != w.ch {
// 			t.Errorf("cells[%d] = {%d,%d}, want {%d,%d}",
// 				i, cells[i].Px, cells[i].Ch, w.px, w.ch)
// 		}
// 		if cells[i].Bits != 2 {
// 			t.Errorf("cells[%d].Bits = %d, want 2", i, cells[i].Bits)
// 		}
// 	}
// }

// func Test_buildCells_SingleChan(t *testing.T) {
// 	chans := []img.ColorChan{img.Red}
// 	cells := buildCells(3, 5, chans)
// 	for i, c := range cells {
// 		if c.Px != i || c.Ch != img.Red || c.Bits != 5 {
// 			t.Errorf("cells[%d] = {%d,%d,%d}, want {%d,Red,5}", i, c.Px, c.Ch, c.Bits, i)
// 		}
// 	}
// }

// // =====================================================================
// // Options.go
// // =====================================================================

// func Test_NewOptions(t *testing.T) {
// 	if _, err := NewOptions(1, []img.ColorChan{img.Red}); err != nil {
// 		t.Errorf("NewOptions(1,[Red]) error: %v", err)
// 	}
// 	if _, err := NewOptions(0, []img.ColorChan{img.Red}); err == nil {
// 		t.Error("NewOptions(0,..) expected error")
// 	}
// 	if _, err := NewOptions(9, []img.ColorChan{img.Red}); err == nil {
// 		t.Error("NewOptions(9,..) expected error")
// 	}
// 	if _, err := NewOptions(1, []img.ColorChan{img.Red, img.Green, img.Blue, img.Red}); err == nil {
// 		t.Error("NewOptions with 4 chans expected error")
// 	}
// }

// func Test_Options_Validate(t *testing.T) {
// 	cases := []struct {
// 		name    string
// 		opt     Options
// 		wantErr bool
// 	}{
// 		{"valid-min", NewOptionsMinimal(), false},
// 		{"valid-opt", NewOptionsOptimal(), false},
// 		{"valid-red", NewOptionsRedChan(), false},
// 		{"valid-agr", NewOptionsAgressive(), false},
// 		{"no-chans", Options{BitsToMod: 1, ChansToEncode: nil}, true},
// 		{"bits-0", Options{BitsToMod: 0, ChansToEncode: []img.ColorChan{img.Red}}, true},
// 		{"bits-9", Options{BitsToMod: 9, ChansToEncode: []img.ColorChan{img.Red}}, true},
// 	}
// 	for _, c := range cases {
// 		err := c.opt.Validate()
// 		if (err != nil) != c.wantErr {
// 			t.Errorf("%s: Validate() err = %v, wantErr = %v", c.name, err, c.wantErr)
// 		}
// 	}
// }

// // =====================================================================
// // Encoder / Decoder
// // =====================================================================

// func Test_NewEncoder_BadOptions(t *testing.T) {
// 	bmp := img.NewBitmap(4, 4)
// 	_, err := NewEncoder(bmp, Options{BitsToMod: 0, ChansToEncode: []img.ColorChan{img.Red}})
// 	if err == nil {
// 		t.Error("NewEncoder: expected error for bad options")
// 	}
// }

// func Test_NewEncoder_BadImage(t *testing.T) {
// 	_, err := NewEncoder(img.Bitmap{}, NewOptionsMinimal())
// 	if err == nil {
// 		t.Error("NewEncoder: expected error for empty image")
// 	}
// }

// func Test_NewDecoder_BadOptions(t *testing.T) {
// 	bmp := img.NewBitmap(4, 4)
// 	_, err := NewDecoder(bmp, Options{BitsToMod: 0, ChansToEncode: []img.ColorChan{img.Red}})
// 	if err == nil {
// 		t.Error("NewDecoder: expected error for bad options")
// 	}
// }

// func Test_NewDecoder_BadImage(t *testing.T) {
// 	_, err := NewDecoder(img.Bitmap{}, NewOptionsMinimal())
// 	if err == nil {
// 		t.Error("NewDecoder: expected error for empty image")
// 	}
// }

// // Round-trip с рандомными данными на всех конфигурациях
// func Test_EncodeDecode_RoundTrip_AllConfigs(t *testing.T) {
// 	rng := rand.New(rand.NewSource(42))

// 	for _, cfg := range configs() {
// 		cfg := cfg
// 		t.Run(cfg.name, func(t *testing.T) {
// 			// достаточно большой образ, чтобы поместились разные длины
// 			baseImg := img.NewBitmap(64, 64)

// 			enc, err := NewEncoder(baseImg, cfg.opt)
// 			if err != nil {
// 				t.Fatalf("NewEncoder: %v", err)
// 			}
// 			dec, err := NewDecoder(baseImg, cfg.opt)
// 			if err != nil {
// 				t.Fatalf("NewDecoder: %v", err)
// 			}

// 			maxBytes := MaxBitsToMode(baseImg, cfg.opt)/8 - 4
// 			if maxBytes <= 0 {
// 				t.Fatalf("bad capacity for %s", cfg.name)
// 			}

// 			for _, n := range []int{0, 1, 2, 3, 7, 16, 64, maxBytes / 2, maxBytes} {
// 				if n < 0 || n > maxBytes {
// 					continue
// 				}
// 				sample := make([]byte, n)
// 				rng.Read(sample)

// 				if err := enc.EncodeBytes(sample); err != nil {
// 					t.Fatalf("%s: EncodeBytes(len=%d) error: %v", cfg.name, n, err)
// 				}
// 				got, err := dec.DecodeBytes(enc.GetEncoded())
// 				if err != nil {
// 					t.Fatalf("%s: DecodeBytes(len=%d) error: %v", cfg.name, n, err)
// 				}
// 				if !bytes.Equal(got, sample) {
// 					t.Fatalf("%s: roundtrip mismatch len=%d\n got=%v\nwant=%v",
// 						cfg.name, n, got, sample)
// 				}
// 			}
// 		})
// 	}
// }

// // EncodeBytes должен возвращать ошибку при превышении ёмкости
// func Test_EncodeBytes_Oversized(t *testing.T) {
// 	for _, cfg := range configs() {
// 		cfg := cfg
// 		t.Run(cfg.name, func(t *testing.T) {
// 			baseImg := img.NewBitmap(2, 2) // 4 пикселя
// 			maxBits := MaxBitsToMode(baseImg, cfg.opt)
// 			// +1 байт сверх ёмкости
// 			tooBig := make([]byte, maxBits/8+1)

// 			enc, err := NewEncoder(baseImg, cfg.opt)
// 			if err != nil {
// 				t.Fatalf("NewEncoder: %v", err)
// 			}
// 			if err := enc.EncodeBytes(tooBig); err == nil {
// 				t.Errorf("expected oversized error")
// 			}
// 		})
// 	}
// }

// // Повторный вызов EncodeBytes перезаписывает сообщение, а не добавляет
// func Test_EncodeBytes_Overwrite(t *testing.T) {
// 	opt := NewOptionsOptimal()
// 	baseImg := img.NewBitmap(32, 32)

// 	enc, _ := NewEncoder(baseImg, opt)
// 	dec, _ := NewDecoder(baseImg, opt)

// 	if err := enc.EncodeBytes([]byte("first message")); err != nil {
// 		t.Fatalf("encode 1: %v", err)
// 	}
// 	if err := enc.EncodeBytes([]byte("second")); err != nil {
// 		t.Fatalf("encode 2: %v", err)
// 	}

// 	got, err := dec.DecodeBytes(enc.GetEncoded())
// 	if err != nil {
// 		t.Fatalf("decode: %v", err)
// 	}
// 	if string(got) != "second" {
// 		t.Errorf("got %q, want %q", got, "second")
// 	}
// }

// // EncodeBytes с пустыми данными: в изображении должен остаться только префикс
// func Test_EncodeDecode_EmptyData(t *testing.T) {
// 	opt := NewOptionsOptimal()
// 	baseImg := img.NewBitmap(16, 16)

// 	enc, _ := NewEncoder(baseImg, opt)
// 	dec, _ := NewDecoder(baseImg, opt)

// 	if err := enc.EncodeBytes([]byte{}); err != nil {
// 		t.Fatalf("EncodeBytes(empty): %v", err)
// 	}
// 	got, err := dec.DecodeBytes(enc.GetEncoded())
// 	if err != nil {
// 		t.Fatalf("DecodeBytes: %v", err)
// 	}
// 	if len(got) != 0 {
// 		t.Errorf("expected empty result, got %v", got)
// 	}
// }

// // DecodeBytes: пустое изображение
// func Test_DecodeBytes_EmptyImage(t *testing.T) {
// 	baseImg := img.NewBitmap(4, 4)
// 	dec, _ := NewDecoder(baseImg, NewOptionsMinimal())

// 	if _, err := dec.DecodeBytes(img.Bitmap{}); err == nil {
// 		t.Error("expected error for empty encoded image")
// 	}
// }

// // DecodeBytes: несовпадение размеров
// func Test_DecodeBytes_SizeMismatch(t *testing.T) {
// 	baseImg := img.NewBitmap(4, 4)
// 	other := img.NewBitmap(8, 8)
// 	dec, _ := NewDecoder(baseImg, NewOptionsMinimal())

// 	if _, err := dec.DecodeBytes(other); err == nil {
// 		t.Error("expected size mismatch error")
// 	}
// }

// // DecodeBytes: слишком маленький образ
// func Test_DecodeBytes_TooSmall(t *testing.T) {
// 	// 1 пиксель, 1 канал, 1 бит — заведомо мало для 4-байтового префикса
// 	baseImg := img.NewBitmap(1, 1)
// 	dec, _ := NewDecoder(baseImg, NewOptionsMinimal())

// 	if _, err := dec.DecodeBytes(baseImg); err == nil {
// 		t.Error("expected 'image too small' error")
// 	}
// }

// // DecodeBytes: повреждённый префикс (слишком большая длина)
// func Test_DecodeBytes_CorruptedLength(t *testing.T) {
// 	opt := NewOptionsOptimal()
// 	baseImg := img.NewBitmap(16, 16)

// 	dec, _ := NewDecoder(baseImg, opt)

// 	// Сконструируем изображение, где в первые 4 байта записана огромная длина.
// 	// Проще всего — закодировать валидное сообщение, а потом вручную поменять
// 	// биты префикса, записав 0xFFFFFFFF.
// 	enc, _ := NewEncoder(baseImg, opt)
// 	if err := enc.EncodeBytes([]byte("hi")); err != nil {
// 		t.Fatalf("encode: %v", err)
// 	}
// 	encImg := enc.GetEncoded()

// 	// Запишем 0xFFFFFFFF в первые 32 бита через writeNBits
// 	bad := encImg.Clone()
// 	prefix := []byte{0xFF, 0xFF, 0xFF, 0xFF}
// 	// используем те же ячейки, что и кодер
// 	cells := buildCells(baseImg.LenPx(), opt.BitsToMod, opt.ChansToEncode)
// 	ibit := 0
// 	for _, cl := range cells {
// 		if ibit >= 32 {
// 			break
// 		}
// 		nextBits := readNBits(prefix, ibit, cl.Bits)
// 		orgByte := baseImg.Px[cl.Px][cl.Ch]
// 		encByte := encodeNBits(orgByte, nextBits, cl.Bits)
// 		bad.Px[cl.Px][cl.Ch] = encByte
// 		ibit += cl.Bits
// 	}

// 	if _, err := dec.DecodeBytes(bad); err == nil {
// 		t.Error("expected error for corrupted length")
// 	}
// }

// // io.EOF: если по какой-то причине битов не хватит (сейчас это ловится раньше,
// // но проверим, что функция не паникует на маленьких изображениях)
// func Test_EncodeBytes_NoPanicOnTiny(t *testing.T) {
// 	baseImg := img.NewBitmap(1, 1)
// 	enc, err := NewEncoder(baseImg, NewOptionsMinimal())
// 	if err != nil {
// 		t.Fatalf("NewEncoder: %v", err)
// 	}
// 	// 1 пиксель * 1 бит = 1 бит ёмкости — точно мало
// 	err = enc.EncodeBytes([]byte("x"))
// 	if err == nil {
// 		t.Error("expected error for tiny image")
// 	}
// 	if errors.Is(err, io.EOF) {
// 		// допустимо
// 	}
// }

// // GetEncoded возвращает копию — изменения не должны влиять на энкодер
// func Test_GetEncoded_ReturnsCopy(t *testing.T) {
// 	opt := NewOptionsOptimal()
// 	baseImg := img.NewBitmap(16, 16)
// 	enc, _ := NewEncoder(baseImg, opt)
// 	_ = enc.EncodeBytes([]byte("abc"))

// 	a := enc.GetEncoded()
// 	b := enc.GetEncoded()

// 	if len(a.Px) != len(b.Px) {
// 		t.Fatal("length mismatch")
// 	}
// 	a.Px[0][0] ^= 0xFF
// 	if a.Px[0][0] == b.Px[0][0] {
// 		t.Error("GetEncoded should return independent copies")
// 	}

// 	c := enc.GetEncoded()
// 	if c.Px[0][0] != b.Px[0][0] {
// 		t.Error("previous copy mutation affected encoder state")
// 	}
// }

// // Изображение-источник не должно меняться при кодировании
// func Test_EncodeBytes_DoesNotMutateBase(t *testing.T) {
// 	opt := NewOptionsOptimal()
// 	baseImg := img.NewBitmap(16, 16)
// 	// заполним детерминированно
// 	for i := range baseImg.Px {
// 		baseImg.Px[i] = img.Pixel{byte(i), byte(i * 2), byte(i * 3)}
// 	}
// 	snapshot := baseImg.Clone()

// 	enc, _ := NewEncoder(baseImg, opt)
// 	if err := enc.EncodeBytes([]byte("hello world")); err != nil {
// 		t.Fatalf("EncodeBytes: %v", err)
// 	}

// 	for i := range baseImg.Px {
// 		if baseImg.Px[i] != snapshot.Px[i] {
// 			t.Fatalf("base image mutated at %d: %v != %v",
// 				i, baseImg.Px[i], snapshot.Px[i])
// 		}
// 	}
// }
