package types

import "testing"

func TestQuanityGet(t *testing.T) {
	cases := []struct {
		name    string
		in      Quanity
		want    int
		wantErr bool
	}{
		{"positive", Quanity(5), 5, false},
		{"one", Quanity(1), 1, false},
		{"zero", QuanityZero, int(QuanityZero), true},
		{"negative", Quanity(-3), int(QuanityInvalid), true},
		{"invalid constant", QuanityInvalid, int(QuanityInvalid), true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := c.in.Get()
			if (err != nil) != c.wantErr {
				t.Fatalf("Get() error = %v, wantErr = %v", err, c.wantErr)
			}
			if got != c.want {
				t.Errorf("Get() = %d, want %d", got, c.want)
			}
		})
	}
}

func TestQuanityAdd(t *testing.T) {
	cases := []struct {
		name string
		a, b Quanity
		want Quanity
	}{
		{"positive + positive", Quanity(3), Quanity(4), Quanity(7)},
		{"one + one", Quanity(1), Quanity(1), Quanity(2)},
		{"zero + positive", QuanityZero, Quanity(5), Quanity(5)},
		{"positive + zero", Quanity(5), QuanityZero, Quanity(5)},
		{"zero + zero", QuanityZero, QuanityZero, QuanityZero},
		{"positive + negative", Quanity(5), Quanity(-1), QuanityInvalid},
		{"negative + positive", Quanity(-1), Quanity(5), QuanityInvalid},
		{"invalid + positive", QuanityInvalid, Quanity(5), QuanityInvalid},
		{"positive + invalid", Quanity(5), QuanityInvalid, QuanityInvalid},
		{"invalid + invalid", QuanityInvalid, QuanityInvalid, QuanityInvalid},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.a.Add(c.b); got != c.want {
				t.Errorf("%v.Add(%v) = %v, want %v", c.a, c.b, got, c.want)
			}
		})
	}
}

func TestQuanityAddIsCommutative(t *testing.T) {
	values := []Quanity{QuanityZero, Quanity(1), Quanity(10), QuanityInvalid, Quanity(-5)}
	for _, a := range values {
		for _, b := range values {
			ab := a.Add(b)
			ba := b.Add(a)
			if ab != ba {
				t.Errorf("Add not commutative: %v+%v=%v, %v+%v=%v",
					a, b, ab, b, a, ba)
			}
		}
	}
}
