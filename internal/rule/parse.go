package rule

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Parser struct {
	service *Service
}

func NewParser(service *Service) *Parser {
	return &Parser{service: service}
}

// dummy -> 0
// dummy -> 1 dummy of MIN,MAX
// dummy -> 2 dummy of AND,OR,NOT

func LastLetter(word string) rune {
	return rune(word[len(word)-1])
}
func ValidOp(op string) bool {
	ops := []string{"<=", "<", "=", "!=", ">", ">="}
	for _, v := range ops {
		if v == op {
			return true
		}
	}
	fmt.Println("INVALID OPERATOR: ", op)
	return false
}
func ValidNum(num string) int {
	intVal, err := strconv.Atoi(num)
	if err == nil {
		return intVal
	}
	fmt.Println("INVALID NUMBER: ", num)
	return -1
}
func (p *Parser) ParseCount(cond string, dummy int) (string, string, int) {
	// considerinf format Count "----"
	// return Matcher, Ineq, Limit
	l := utf8.RuneCountInString(cond)
	words := strings.Fields(cond)
	len1 := len(words)

	Limit := -1
	if dummy != 1 {
		Limit = ValidNum(words[len1-1])
	}
	if dummy == 1 {
		if words[1][0] != '"' || LastLetter(words[len1-1]) != '"' {
			println("InValid representation for Count Instruction in Min/Max.")
			println("Please follow: Count \"your_string\", don't provide operator or number")
			return "", "", -1
		}
	} else {
		if words[1][0] != '"' || LastLetter(words[len1-3]) != '"' || !ValidOp(words[len1-2]) || Limit == -1 {
			println("InValid representation for Count Instruction.")
			println("Please follow: Count \"your_string\" operator number")
			return "", "", -1
		}
	}

	if dummy == 1 {
		return cond[7 : l-1], "", 0
	}

	return cond[7 : l-len(words[len1-1])-len(words[len1-2])-3], words[len1-2], Limit
}
func (p *Parser) ParseLength(cond string, dummy int) (string, int) {
	words := strings.Fields(cond)
	len1 := len(words)
	Limit := -1
	if len1 == 3 {
		Limit = ValidNum(words[len1-1])
	}
	if dummy == 1 {
		if len1 != 1 {
			println("InValid representation for Length Instruction in Min/Max.")
			println("Please Just mention: Length")
			return "", -1
		}
	} else {
		if len1 != 3 || !ValidOp(words[1]) || Limit == -1 {
			println("InValid representation for Length Instruction.")
			println("Please follow: Length operator number")
			return "", -1
		}
	}
	if dummy == 1 {
		return "", 0
	}
	return words[1], Limit
}
func (p *Parser) ParseContains(cond string) (string, int) {
	l := len(cond)
	if cond[9] != '"' || cond[l-1] != '"' {
		println("InValid representation for Contains Instruction.")
		println("Please follow: Contains \"your_string\"")
		return "", -1
	}
	return cond[10 : l-1], 0
}
func (p *Parser) ParseMinMax(ctx context.Context, cond string, dummy int) ([]primitive.ObjectID, string, int) {
	words := strings.Fields(cond)
	len1 := len(words)
	obj := []primitive.ObjectID{}
	var Nobj primitive.ObjectID
	Limit := -1
	if dummy != 1 {
		Limit = ValidNum(words[len1-1])
	}
	if dummy == 1 {
		if LastLetter(cond) != ')' || cond[4] != '(' {
			println("InValid representation for Min/Max Instruction inside Min/Max.")
			println("Please follow: MIN/MAX (Instruction1 , Instruction2 , ... , InstructionN)")
			return obj, "", -1
		}
	} else {
		if cond[4] != '(' || LastLetter(words[len1-3]) != ')' || !ValidOp(words[len1-2]) || Limit == -1 {
			println("InValid representation for Min/Max Instruction.")
			println("Please follow: MIN/MAX (Instruction1 , Instruction2 , ... , InstructionN) operator number")
			return obj, "", -1
		}
	}
	// ip for parenthesis balance count
	// ib for brackects balance count
	var ip, ib, cond_start int
	cond_start = 5
	for i, v := range cond {
		if v == '(' {
			ip++
		} else if v == ')' {
			ip--
		}
		if v == '"' {
			ib = (ib + 1) % 2
		}
		if v == ',' && ip == 1 && ib == 0 {
			if cond[i-1] != ' ' || cond[i+1] != ' ' {
				fmt.Println("Please use spaces and commas',' between Instructions in Min/Max Eg:- MIN(COUNT \"a\" , COUNT \"b\") > 4")
				return []primitive.ObjectID{}, "", -1
			}
			if cond[i-2] == ' ' || cond[i+2] == ' ' {
				fmt.Println("Please use only one single spaces between commas',' and Instructions in Min/Max Eg:- MIN(COUNT \"a\" , COUNT \"b\") > 4")
				return []primitive.ObjectID{}, "", -1
			}
		}
		if (v == ',' && ip == 1 && ib == 0) || (v == ')' && ip == 0) {
			// this code checks if v == ')' this is last parenthesis after this only op number are their
			if v == ')' && (i != len(cond)-1 && (i+2+len(words[len1-2])+len(words[len1-1]) != len(cond)-1)) {
				println("InValid representation for Min/Max Instruction.")
				println("After Closing Parenthesis of MIN/MAX Instruction their should be only operator and number")
				println("Please follow: MIN/MAX (Instruction1 , Instruction2 , ... , InstructionN) operator number")
				return obj, "", -1
			}
			if v == ',' {
				Nobj = p.ParseCondition(ctx, cond[cond_start:i-1], 1, "")
				cond_start = i + 2
			} else {
				Nobj = p.ParseCondition(ctx, cond[cond_start:i], 1, "")
			}
			if Nobj == primitive.NilObjectID {
				return []primitive.ObjectID{}, "", -1
			}
			obj = append(obj, Nobj)
		}
	}
	if dummy == 1 {
		return obj, "", 0
	}
	return obj, words[len1-2], Limit
}
func (p *Parser) ParseAndOr(ctx context.Context, cond string) ([]primitive.ObjectID, int) {
	// PLEASE ONCE GO THRU THIS CODE

	var obj = []primitive.ObjectID{}
	var Nobj primitive.ObjectID
	// ip for parenthesis balance count
	// ib for brackects balance count
	var ii, ip, ib, cond_start int
	if cond[:3] == "AND" {
		ii = 4
		cond_start = 5
	} else {
		ii = 3
		cond_start = 4
	}
	if LastLetter(cond) != ')' || cond[ii] != '(' {
		println("InValid representation for Min/Max Instruction inside Min/Max.")
		println("Please follow: MIN/MAX (Instruction1 , Instruction2 , ... , InstructionN)")
		return obj, -1
	}

	for i, v := range cond {
		if v == '(' {
			ip++
		} else if v == ')' {
			ip--
		}
		if v == '"' {
			ib = (ib + 1) % 2
		}
		if v == ',' && ip == 1 && ib == 0 {
			if cond[i-1] != ' ' || cond[i+1] != ' ' {
				fmt.Println("Please use spaces and commas',' between Instructions in AND/OR Eg:- AND(COUNT \"a\"  < 4 , COUNT \"b\" > 3)")
				return []primitive.ObjectID{}, -1
			}
			if cond[i-2] == ' ' || cond[i+2] == ' ' {
				fmt.Println("Please use only one single spaces between commas',' and Instructions in AND/OR Eg:- AND(COUNT \"a\"  < 4 , COUNT \"b\" > 3)")
				return []primitive.ObjectID{}, -1
			}
		}
		if (v == ',' && ip == 1 && ib == 0) || (v == ')' && ip == 0) {
			// this code checks if v == ')' this is last parenthesis after this only op number are their
			if v == ')' && i != len(cond)-1 {
				println("InValid representation for AND/OR Instruction.")
				println("After Closing Parenthesis of AND/OR Instruction their should not be operator and number")
				println("Please follow: AND/OR (Instruction1 , Instruction2 , ... , InstructionN)")
				return obj, -1
			}
			if v == ',' {
				Nobj = p.ParseCondition(ctx, cond[cond_start:i-1], 2, "")
				cond_start = i + 2
			} else {
				Nobj = p.ParseCondition(ctx, cond[cond_start:i], 2, "")
			}
			if Nobj == primitive.NilObjectID {
				return []primitive.ObjectID{}, -1
			}
			obj = append(obj, Nobj)
		}
	}
	return obj, 0
}
func (p *Parser) ParseNot(ctx context.Context, cond string) ([]primitive.ObjectID, int) {
	l := len(cond)
	var obj = []primitive.ObjectID{}
	if cond[4] != '(' || cond[l-1] != ')' {
		println("InValid representation for NOT Instruction.")
		println("Please follow: NOT (Instruction")
	}
	Nobj := p.ParseCondition(ctx, cond[5:l-1], 2, "")
	if Nobj == primitive.NilObjectID {
		return obj, -1
	}
	obj = append(obj, Nobj)
	return obj, 0
}

var owner primitive.ObjectID

func (p *Parser) ParseCondition(ctx context.Context, cond string, dummy int, temp string) primitive.ObjectID {
	//fmt.Println("\n\n\n",cond,"\n")
	var ip, ib int
	for _, v := range cond {
		if v == '(' {
			ip++
		} else if v == ')' {
			ip--
		}
		if v == '"' {
			ib = (ib + 1) % 2
		}
		if v == ',' && ip == 0 && ib == 0 {
			fmt.Println("If you are using Multiple Instructions using comma\",\" Please use one of AND, OR, MAX, MIN")
			return primitive.NilObjectID
		}
	}
	Names := []string{"Count", "Length", "Contains", "MAX", "MIN", "OR", "AND", "NOT"}
	//var ii int
	var tag int
	var Name, Matcher, Ineq, Notify, When string
	//var IsTheir bool
	var Limit int
	var Obj = []primitive.ObjectID{}
	//if dummy == 0 {
	Notify = temp
	When = cond
	//}
	for i, v := range cond {
		//println(i,v,' ','(')
		if v == '(' || v == '"' {
			fmt.Println("Please provide space between Instruction name and ", string(v), ".")
			if v == '(' {
				fmt.Println(" Example: MIN (condition1 , condition2) ")
			} else {
				fmt.Println(" Example: Count \"abc\" operater number")
			}
			return primitive.NilObjectID
		}
		if v == ' ' || i == len(cond)-1 {
			if i == len(cond)-1 && v != ' ' {
				Name = cond[:i+1]
			} else {
				Name = cond[:i]
			}
			//ii = i + 1
			var have bool
			for _, vv := range Names {
				if vv == Name {
					have = true
					break
				}
			}
			if !have {
				fmt.Println("Not supports your condition", Name)
				return primitive.NilObjectID
			}
			break
		}
		if i == 8 {
			fmt.Println("Not supports your condition. Supports only: ", Names)
			return primitive.NilObjectID
		}
	}
	//println(ii)
	if dummy == 1 {
		if Name == "Contains" || Name == "AND" || Name == "OR" || Name == "NOT" {
			fmt.Println("You cant have ", Name, "inside Min/Max Instruction")
			return primitive.NilObjectID
		}
	}
	if Name == "Count" {
		fmt.Println("Gointing into ParseCount")
		Matcher, Ineq, Limit = p.ParseCount(cond, dummy)
		fmt.Println("Count: ", Matcher, ',', Ineq, ',', Limit)
		if Limit == -1 {
			return primitive.NilObjectID
		}
	} else if Name == "Length" {
		fmt.Println("Going into ParseLength")
		Ineq, Limit = p.ParseLength(cond, dummy)
		fmt.Println("Length: ", Ineq, ',', Limit)
		if Limit == -1 {
			return primitive.NilObjectID
		}
	} else if Name == "Contains" {
		fmt.Println("Going into ParseContains")
		Matcher, tag = p.ParseContains(cond)
		fmt.Println("Contains: ", Matcher)
		if tag == -1 {
			return primitive.NilObjectID
		}
	} else if Name == "MAX" || Name == "MIN" {
		fmt.Println("Going into ParseMinMax")
		Obj, Ineq, Limit = p.ParseMinMax(ctx, cond, dummy)
		/*fmt.Println(Name, ": ")
		for i, v := range Obj {
			fmt.Println(i, v)
		}*/
		fmt.Println(Name, Ineq, ',', Limit)
		if Limit == -1 {
			return primitive.NilObjectID
		}
	} else if Name == "AND" || Name == "OR" {
		fmt.Println("Going into ParseAndOr")
		Obj, tag = p.ParseAndOr(ctx, cond)
		if tag == -1 {
			return primitive.NilObjectID
		}
	} else if Name == "NOT" {
		fmt.Println("Going into ParseNot")
		Obj, tag = p.ParseNot(ctx, cond)
		if tag == -1 {
			return primitive.NilObjectID
		}
	} else {
		fmt.Println("InValid Instruction: ", Name)
		return primitive.NilObjectID
	}
	var Nrule Rule
	Nrule.Name = Name
	Nrule.Matcher = Matcher
	Nrule.Ineq = Ineq
	Nrule.Limit = Limit
	Nrule.Obj = Obj
	Nrule.Notify = Notify
	Nrule.When = When
	Nrule.Owners = []primitive.ObjectID{}

	var ID primitive.ObjectID
	var err error
	if dummy == 0 {
		ID, err = p.service.FindDoc(ctx, &Nrule, owner)
	} else {
		ID, err = p.service.FindDoc(ctx, &Nrule, primitive.NilObjectID)
	}

	if err != nil {
		if dummy == 0 {
			Nrule.Owners = append(Nrule.Owners, owner)
			ID, err = p.service.CreateRule(ctx, &Nrule, owner)
		} else {
			ID, err = p.service.CreateRule(ctx, &Nrule, primitive.NilObjectID)
		}
		if err != nil {
			println(err.Error())
			return primitive.NilObjectID
		}
	}
	// fmt.Println("printing after ParseCondition:\n", Name, Matcher, Ineq, Notify, When, Obj, Limit, ID, Created)
	return ID
}

func (p *Parser) ParseRule(ctx context.Context, text string, rule_owner primitive.ObjectID) (Rule, error) {
	var rule Rule
	notifyPattern := regexp.MustCompile(`NOTIFY\s+(.*?)\s+WHEN\s+(.*)`)
	matches := notifyPattern.FindStringSubmatch(text)
	if len(matches) != 3 {
		return rule, errors.New("not of form: NOTIFY \"your_classification\" WHEN \"your_condition\"")
	}

	owner = rule_owner
	ID := p.ParseCondition(ctx, matches[2], 0, matches[1])

	// What happens when ID is nilObjectID
	err := p.service.repo.collection.FindOne(ctx, bson.D{{Key: "_id", Value: ID}}).Decode(&rule)

	return rule, err
}
