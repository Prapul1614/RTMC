package rule

import (
	"context"
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

func TestGetContains(t *testing.T) {
	if !GetContains("barbabarber", "babar") {
		t.Fatal("Unable to find if pattern is contained")
	}

	if GetContains("yreuyberublinbf", "a") {
		t.Fatal("Detected patteren but actually not present")
	}
}

func TestImplementOp(t *testing.T) {
	if !ImplementOp("<", 4, 5) {
		t.Fatal("Saying 4 >= 5")
	}

	if !ImplementOp("=", 5, 5) {
		t.Fatal("Saying 5 != 5")
	}

	if !ImplementOp("!=", 4, 5) {
		t.Fatal("Saying 4 == 5")
	}

	if !ImplementOp("<=", 4, 4) {
		t.Fatal("Saying 4 > 4")
	}

	if !ImplementOp(">=", 5, 4) {
		t.Fatal("Saying 5 < 4")
	}

	if !ImplementOp(">", 10, 2) {
		t.Fatal("Saying 10 <= 2")
	}
}


func TestImplementMinMax(t *testing.T) {
	testservice := &Service{mockrepo: &MockRepository{}}
    ctx := context.Background()
	
    // Test cases
    tests := []struct {
        text     string
        rule     *Rule
        expected int
    }{
        {"hello", mockRules[ruleIDs[0]], 5},
        {"aa ahahbahbavuta", mockRules[ruleIDs[1]], 7},
        {"hello aa", mockRules[ruleIDs[2]], 2},
        {"hbnnjbhgb", mockRules[ruleIDs[3]], 3},
		{"jnabaybbaybayubdbdabbaab", mockRules[ruleIDs[4]], 9},
    }

    for i, tc := range tests {
		if i >= 5 {break}
		result := testservice.ImplementMinMax(ctx, tc.text, tc.rule)
		if result != tc.expected {
			t.Errorf("expected %d, got %d. for %v - %v", tc.expected, result, i , mockRules[ruleIDs[i]])
		}
    }
}

func TestImplementAndOr(t *testing.T) {
	testservice := &Service{mockrepo: &MockRepository{}}
    ctx := context.Background()
    // Test cases
    tests := []struct {
        text     string
        rule     *Rule
        expected bool
    }{
        {"hell0", mockRules[ruleIDs[7]], false},
        {"aa ahahbahbavuta", mockRules[ruleIDs[7]], true},
        {"haaaello aa", mockRules[ruleIDs[9]], false},
        {"habb", mockRules[ruleIDs[9]], true},
		{"jnabayhuhyuhhh", mockRules[ruleIDs[9]], true},
    }

    for i, tc := range tests {
		result := testservice.ImplementAndOr(ctx, tc.text, tc.rule)
		if result != tc.expected {
			t.Errorf("expected %v, got %v. for %v - %v", tc.expected, result, i , mockRules[ruleIDs[i]])
		}
    }
}

