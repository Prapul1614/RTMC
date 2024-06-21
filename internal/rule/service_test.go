package rule

import (
	"math/rand"
	"testing"
)

func TestGetLength(t *testing.T) {
	l := rand.Intn(100)
	str := randomString(l)

	if GetLength(str) != l {
		t.Fatalf("Giving wrong length")
	}
}

func TestGetCount(t *testing.T) {
	if GetCount("hihoi hihoi", "hi", false) != 2 {
		t.Fatal("Wrong Count")
	}

	if GetCount("hhhi hhhi", "hh", false) != 4 {
		t.Fatal("Wrong Count")
	}

	if GetCount("hhhi hhhi", "hi", false) != 2 {
		t.Fatal("Wrong Count")
	}
}