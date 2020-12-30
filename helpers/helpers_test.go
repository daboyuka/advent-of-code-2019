package helpers

import "testing"

func TestSplit2(t *testing.T) {
	tests := []struct {
		S, Sep           string
		ExpectA, ExpectB string
	}{
		{"aaabbccccc", "bb", "aaa", "ccccc"},
		{"aaabbccccc", "b", "aaa", "bccccc"},
		{"aaabbccccc", "c", "aaabb", "cccc"},
		{"aaabbccccc", "d", "aaabbccccc", ""},
	}
	for _, tst := range tests {
		a, b := Split2(tst.S, tst.Sep)
		if a != tst.ExpectA || b != tst.ExpectB {
			t.Fatalf("%+v vs. '%s' '%s'", tst, a, b)
		}
	}
}

func TestSplit3(t *testing.T) {
	tests := []struct {
		S, Sep1, Sep2             string
		ExpectA, ExpectB, ExpectC string
	}{
		{"aaabbcccccdddeeee", "bb", "ddd", "aaa", "ccccc", "eeee"},
		{"aaabbcccccdddeeee", "b", "d", "aaa", "bccccc", "ddeeee"},
		{"aaabbcccccdddeeee", "f", "ddd", "aaabbccccc", "", "eeee"},
		{"aaabbcccccdddeeee", "bb", "f", "aaa", "cccccdddeeee", ""},
		{"aaabbcccccdddeeee", "f", "g", "aaabbcccccdddeeee", "", ""},
	}
	for _, tst := range tests {
		a, b, c := Split3(tst.S, tst.Sep1, tst.Sep2)
		if a != tst.ExpectA || b != tst.ExpectB || c != tst.ExpectC {
			t.Fatalf("%+v vs. '%s' '%s' '%s'", tst, a, b, c)
		}
	}
}
