package img

import (
	"image"
	"image/draw"
)

type Bitmap struct {
	Px []Pixel
	W  int
	H  int
}

func NewBitmap(w, h int) Bitmap {
	return Bitmap{
		Px: make([]Pixel, w*h),
		W:  w,
		H:  h,
	}
}

func (b Bitmap) Clone() Bitmap {
	clone := Bitmap{
		Px: make([]Pixel, len(b.Px)),
		W:  b.W,
		H:  b.H,
	}
	copy(clone.Px, b.Px)
	return clone
}

func (b Bitmap) LenPx() int {
	return len(b.Px)
}

// At возвращает пиксель по координатам (x, y).
// Без проверки границ — как в image.Image.
func (b Bitmap) At(x, y int) Pixel {
	return b.Px[y*b.W+x]
}

// Set устанавливает пиксель по координатам (x, y).
func (b Bitmap) Set(x, y int, p Pixel) {
	b.Px[y*b.W+x] = p
}

// --- Конвертация с image.Image ---

// FromImage конвертирует стандартный image.Image в Bitmap.
// Пиксели берутся из NRGBA-представления, чтобы одинаково работать
// с любыми моделями цвета (RGB, RGBA, палитра и т.д.).
func FromImage(src image.Image) Bitmap {
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()

	out := NewBitmap(w, h)

	nrgba := image.NewNRGBA(image.Rect(0, 0, w, h))
	draw.Draw(nrgba, nrgba.Bounds(), src, b.Min, draw.Src)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			i := y*nrgba.Stride + x*4
			out.Px[y*w+x] = Pixel{
				nrgba.Pix[i],   // R
				nrgba.Pix[i+1], // G
				nrgba.Pix[i+2], // B
			}
		}
	}
	return out
}

// ToImage конвертирует Bitmap в image.Image (NRGBA).
func (b Bitmap) ToImage() image.Image {
	out := image.NewNRGBA(image.Rect(0, 0, b.W, b.H))
	for y := 0; y < b.H; y++ {
		for x := 0; x < b.W; x++ {
			px := b.Px[y*b.W+x]
			j := y*out.Stride + x*4
			out.Pix[j] = px[0]
			out.Pix[j+1] = px[1]
			out.Pix[j+2] = px[2]
			out.Pix[j+3] = 255
		}
	}
	return out
}
