package rule

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)
var service = &Service{}
var parser = NewParser(service)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
func randomString (l int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, l)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func TestLastLetter(t *testing.T) {
	if LastLetter("Hell o") != rune('o') {
		t.Fatalf("Incorrect Last letter")
	}
}

func TestValidOp(t *testing.T) {
	op := randomString(rand.Intn(2) + 1)
	var found bool

	ops := []string{"<", "<=", "=", "!=", ">=", ">"}
	for _,v := range(ops) {
		if !ValidOp(v) {
			t.Fatalf("Declaring %v as not a valid operator", op)
		}
		if v == op {found = true}
	}

	if !found && ValidOp(op){
		t.Fatalf("Should have declared \"%v\" not an operator", op)
	}
}

func TestValidNum(t *testing.T) {
	if ValidNum("-2") !=  -1 {
		t.Fatalf("Declaring %v as valid number", "-2")
	}
	if ValidNum("2.3") != -1 {
		t.Fatalf("Declaring %v as valid number", "2.3")
	}
	if ValidNum("2/3") != -1 {
		t.Fatalf("Declaring %v as valid number", "2/3")
	}
	if ValidNum("5") != 5 {
		t.Fatalf("Wasn't able to convert string to currect number")
	}
}

func TestParseCount(t *testing.T) {
	cond := fmt.Sprintf("Count \"%v\" ", "hii")
	matcher, op, limit, msg := parser.ParseCount(cond, 1)
	if msg != "Success" || matcher != "hii" || op != "" || limit != 0{
		t.Fatalf("Error: %v", msg)
	}

	_, _, _, msg = parser.ParseCount(cond, 0)
	if !strings.HasSuffix(msg, "Please follow: Count \"your_string\" operator number"){
		t.Fatalf("Error, Should have given Please follow: Count \"your_string\" operator number")
	}

	cond = fmt.Sprintf("Count \"%v\" %v %v", "hii", "<=", "5")
	matcher, op, limit, msg = parser.ParseCount(cond, 0)
	if matcher != "hii" || op != "<=" || limit != 5 || msg != "Success" {
		t.Fatalf("Not parsed properly: %v", cond)
	}

	cond = fmt.Sprintf("Count \"%v\" %v %v", "hii", "<=", "5")
	_, _, _, msg = parser.ParseCount(cond, 1)
	if !strings.HasSuffix(msg, "Please follow: Count \"your_string\", don't provide operator or number") {
		t.Fatalf("Error, Should have given Please follow: Count \"your_string\"")
	}

}

func TestParseLength(t *testing.T) {
	cond := "Length < 20"
	op, limit, msg := parser.ParseLength(cond, 0)
	if op != "<" || limit != 20 || msg != "Success" {
		t.Fatalf("Couldn't parse: %v", cond)
	}

	_, _, msg = parser.ParseLength(cond, 1)
	if !strings.HasSuffix(msg, "Please Just mention: Length") {
		t.Fatalf("Error should have been Please Just mention: Length")
	}

	cond = "Length"
	_, _, msg = parser.ParseLength(cond, 0)
	if !strings.HasSuffix(msg, "Please follow: Length operator number") {
		t.Fatalf("Error should have been Please follow: Length operator number")
	}
}

func TestParseContains(t *testing.T) {
	cond := "Contains \"hi\""
	matcher, _, msg := parser.ParseContains(cond)
	if matcher != "hi" || msg != "Success" {
		t.Fatalf("Error unable to parse: %v", cond)
	}
}

func TestParseMinMax(t *testing.T) {
	cond := "MIN (Length , Count \"a\")"
	_, _, _, msg := parser.ParseMinMax(context.Background(),cond, 0)
	if !strings.HasSuffix(msg, "Please follow: MIN/MAX (Instruction1 , Instruction2 , ... , InstructionN) operator number") {
		t.Fatalf("Error should be - Please follow: MIN/MAX (Instruction1 , Instruction2 , ... , InstructionN) operator number")
	}

	cond = "MIN (Length , Count \"a\") != 5"
	_, _, _, msg = parser.ParseMinMax(context.Background(),cond, 1)
	if !strings.HasSuffix(msg, "Please follow: MIN/MAX (Instruction1 , Instruction2 , ... , InstructionN)") {
		t.Fatalf("Error should be - Please follow: MIN/MAX (Instruction1 , Instruction2 , ... , InstructionN)")
	}

	cond = "MIN(Length , Count \"a\") != 5"
	_, _, _, msg = parser.ParseMinMax(context.Background(),cond, 0)
	if !strings.HasSuffix(msg, "Please follow: MIN/MAX (Instruction1 , Instruction2 , ... , InstructionN) operator number") {
		t.Fatalf("Error should be - Please follow: MIN/MAX (Instruction1 , Instruction2 , ... , InstructionN) operator number")
	}
}

func TestParseAndOr(t *testing.T) {
	cond := "AND (Length < 5 , Count \"a\" = 5) != 5"
	_, _, msg := parser.ParseAndOr(context.Background(),cond)
	if !strings.HasSuffix(msg, "Please follow: AND/OR (Instruction1 , Instruction2 , ... , InstructionN)") {
		t.Fatalf("Error should be - Please follow: AND/OR (Instruction1 , Instruction2 , ... , InstructionN)")
	}

	cond = "AND(Length < 5 , Count \"a\" = 5)"
	_, _, msg = parser.ParseAndOr(context.Background(),cond)
	if !strings.HasSuffix(msg, "Please follow: AND/OR (Instruction1 , Instruction2 , ... , InstructionN)") {
		t.Fatalf("Error should be - Please follow: AND/OR (Instruction1 , Instruction2 , ... , InstructionN)")
	}
}