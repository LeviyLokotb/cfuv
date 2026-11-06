package img

import "testing"

func Test_Bitmap_Clone(t *testing.T) {
	bmp := Bitmap{
		Px: []Pixel{
			Pixel{1, 0, 1},
			Pixel{0, 2, 0},
			Pixel{1, 0, 1},
		}}
	clone := bmp.Clone()

	if bmp.LenPx() != clone.LenPx() {
		t.Fatalf("len not match (%d and %d)", bmp.LenPx(), clone.LenPx())
	}

	if &bmp.Px[0] == &clone.Px[0] {
		t.Fatal("clone is not independed")
	}

	for i := range bmp.Px {
		if bmp.Px[i] != clone.Px[i] {
			t.Errorf("pixels not match (%s and %s)", bmp.Px[i], clone.Px[i])
		}
	}
}
