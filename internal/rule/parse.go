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
func (p *Parser) ParseCount(cond string, dummy int) (string, string, int, string) {
	var msg = fmt.Sprintf("Error in formatting of:  %v \n", cond)
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
			msg += "InValid representation for Count Instruction in Min/Max.\n"
			msg += "Please follow: Count \"your_string\", don't provide operator or number"
			return "", "", -1, msg
		}
	} else {
		if words[1][0] != '"' || LastLetter(words[len1-3]) != '"' || !ValidOp(words[len1-2]) || Limit == -1 {
			msg += "InValid representation for Count Instruction.\n"
			msg += "Please follow: Count \"your_string\" operator number"
			return "", "", -1, msg
		}
	}

	if dummy == 1 {
		return cond[7 : l-1], "", 0, "Success"
	}

	return cond[7 : l-len(words[len1-1])-len(words[len1-2])-3], words[len1-2], Limit, "Success"
}
func (p *Parser) ParseLength(cond string, dummy int) (string, int, string) {
	var msg = fmt.Sprintf("Error in formatting of:  %v \n", cond)
	words := strings.Fields(cond)
	len1 := len(words)
	Limit := -1
	if len1 == 3 {
		Limit = ValidNum(words[len1-1])
	}
	if dummy == 1 {
		if len1 != 1 {
			msg += "InValid representation for Length Instruction in Min/Max.\n"
			msg += "Please Just mention: Length"
			return "", -1, msg
		}
	} else {
		if len1 != 3 || !ValidOp(words[1]) || Limit == -1 {
			msg += "InValid representation for Length Instruction.\n"
			msg += "Please follow: Length operator number"
			return "", -1, msg
		}
	}
	if dummy == 1 {
		return "", 0, "Success"
	}
	return words[1], Limit, "Success"
}
func (p *Parser) ParseContains(cond string) (string, int, string) {
	var msg = fmt.Sprintf("Error in formatting of:  %v \n", cond)
	l := len(cond)
	if cond[9] != '"' || cond[l-1] != '"' {
		msg += "InValid representation for Contains Instruction.\n"
		msg += "Please follow: Contains \"your_string\""
		return "", -1, msg
	}
	return cond[10 : l-1], 0, msg
}
func (p *Parser) ParseMinMax(ctx context.Context, cond string, dummy int) ([]primitive.ObjectID, string, int, string) {
	var msg = fmt.Sprintf("Error in formatting of:  %v \n", cond)
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
			msg += "InValid representation for Min/Max Instruction inside Min/Max.\n"
			msg += "Please follow: MIN/MAX (Instruction1 , Instruction2 , ... , InstructionN)"
			return obj, "", -1, msg
		}
	} else {
		if cond[4] != '(' || LastLetter(words[len1-3]) != ')' || !ValidOp(words[len1-2]) || Limit == -1 {
			msg += "InValid representation for Min/Max Instruction.\n"
			msg += "Please follow: MIN/MAX (Instruction1 , Instruction2 , ... , InstructionN) operator number"
			return obj, "", -1, msg
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
				msg += "Please use spaces and commas',' between Instructions in Min/Max Eg:- MIN(COUNT \"a\" , COUNT \"b\") > 4"
				return []primitive.ObjectID{}, "", -1, msg
			}
			if cond[i-2] == ' ' || cond[i+2] == ' ' {
				msg += "Please use only one single spaces between commas',' and Instructions in Min/Max Eg:- MIN(COUNT \"a\" , COUNT \"b\") > 4"
				return []primitive.ObjectID{}, "", -1, msg
			}
		}
		if (v == ',' && ip == 1 && ib == 0) || (v == ')' && ip == 0) {
			if v == ')' && (i != len(cond)-1 && (i+2+len(words[len1-2])+len(words[len1-1]) != len(cond)-1)) {
				msg += "InValid representation for Min/Max Instruction.\n"
				msg += "After Closing Parenthesis of MIN/MAX Instruction their should be only operator and number.\n"
				msg += "Please follow: MIN/MAX (Instruction1 , Instruction2 , ... , InstructionN) operator number."
				return obj, "", -1, msg
			}
			if v == ',' {
				Nobj,msg = p.ParseCondition(ctx, cond[cond_start:i-1], 1, "", primitive.NilObjectID)
				cond_start = i + 2
			} else {
				Nobj,msg = p.ParseCondition(ctx, cond[cond_start:i], 1, "", primitive.NilObjectID)
			}
			if Nobj == primitive.NilObjectID {
				return []primitive.ObjectID{}, "", -1, msg
			}
			obj = append(obj, Nobj)
		}
	}
	if dummy == 1 {
		return obj, "", 0, "Success"
	}
	return obj, words[len1-2], Limit, "Success"
}
func (p *Parser) ParseAndOr(ctx context.Context, cond string) ([]primitive.ObjectID, int, string) {
	
	var msg = fmt.Sprintf("Error in formatting of:  %v \n", cond)
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
		msg += "InValid representation for Min/Max Instruction inside Min/Max.\n"
		msg += "Please follow: MIN/MAX (Instruction1 , Instruction2 , ... , InstructionN)"
		return obj, -1, msg
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
				msg += "Please use spaces and commas',' between Instructions in AND/OR Eg:- AND(COUNT \"a\"  < 4 , COUNT \"b\" > 3)"
				return []primitive.ObjectID{}, -1, msg
			}
			if cond[i-2] == ' ' || cond[i+2] == ' ' {
				msg += "Please use only one single spaces between commas',' and Instructions in AND/OR Eg:- AND(COUNT \"a\"  < 4 , COUNT \"b\" > 3)"
				return []primitive.ObjectID{}, -1, msg
			}
		}
		if (v == ',' && ip == 1 && ib == 0) || (v == ')' && ip == 0) {
			if v == ')' && i != len(cond)-1 {
				msg += "After Closing Parenthesis of AND/OR Instruction their should not be operator and number\n"
				msg += "Please follow: AND/OR (Instruction1 , Instruction2 , ... , InstructionN)"
				return obj, -1, msg
			}
			if v == ',' {
				Nobj,msg = p.ParseCondition(ctx, cond[cond_start:i-1], 2, "", primitive.NilObjectID)
				cond_start = i + 2
			} else {
				Nobj,msg = p.ParseCondition(ctx, cond[cond_start:i], 2, "", primitive.NilObjectID)
			}
			if Nobj == primitive.NilObjectID {
				return []primitive.ObjectID{}, -1, msg
			}
			obj = append(obj, Nobj)
		}
	}
	return obj, 0, "Success"
}
func (p *Parser) ParseNot(ctx context.Context, cond string) ([]primitive.ObjectID, int, string) {
	var msg = fmt.Sprintf("Error in formatting of:  %v \n", cond)
	l := len(cond)
	var obj = []primitive.ObjectID{}
	if cond[4] != '(' || cond[l-1] != ')' {
		msg += "InValid representation for NOT Instruction.\n"
		msg += "Please follow: NOT (Instruction)."
		return obj, -1, msg
	}
	Nobj,msg := p.ParseCondition(ctx, cond[5:l-1], 2, "",primitive.NilObjectID)
	if Nobj == primitive.NilObjectID {
		return obj, -1, msg
	}
	obj = append(obj, Nobj)
	return obj, 0, "Success"
}


func (p *Parser) ParseCondition(ctx context.Context, cond string, dummy int, temp string, owner primitive.ObjectID) (primitive.ObjectID, string) {
	var msg = fmt.Sprintf("Error in formatting of:  %v \n", cond)
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
			msg += "If you are using Multiple Instructions using comma\",\" Please use one of AND, OR, MAX, MIN"
			return primitive.NilObjectID, msg
		}
	}
	Names := []string{"Count", "Length", "Contains", "MAX", "MIN", "OR", "AND", "NOT"}

	var tag int
	var Name, Matcher, Ineq, Notify, When string
	var Limit int
	var Obj = []primitive.ObjectID{}
	
	Notify = temp
	When = cond
	
	for i, v := range cond {
		if v == '(' || v == '"' {
			msg += fmt.Sprintf("Please provide space between Instruction name and %v .\n", string(v))
			if v == '(' {
				msg += " Example: MIN (condition1 , condition2) "
			} else {
				msg += " Example: Count \"abc\" operater number"
			}
			return primitive.NilObjectID, msg
		}
		if v == ' ' || i == len(cond)-1 {
			if i == len(cond)-1 && v != ' ' {
				Name = cond[:i+1]
			} else {
				Name = cond[:i]
			}
			var have bool
			for _, vv := range Names {
				if vv == Name {
					have = true
					break
				}
			}
			if !have {
				msg += fmt.Sprintf("Not supports your condition %v .\n Supports only %v", Name, Names)
				return primitive.NilObjectID, msg
			}
			break
		}
		if i == 8 {
			msg += fmt.Sprintf("Not supports your condition. Supports only %v", Names)
			return primitive.NilObjectID, msg
		}
	}
	
	if dummy == 1 {
		if Name == "Contains" || Name == "AND" || Name == "OR" || Name == "NOT" {
			msg += fmt.Sprintf("You cant have %v inside Min/Max Instruction", Name)
			return primitive.NilObjectID, msg
		}
	}
	
	if Name == "Count" {
		fmt.Println("Gointing into ParseCount")
		Matcher, Ineq, Limit, msg = p.ParseCount(cond, dummy)
		fmt.Println("Count: ", Matcher, ',', Ineq, ',', Limit)
		if Limit == -1 {
			return primitive.NilObjectID, msg
		}
	} else if Name == "Length" {
		fmt.Println("Going into ParseLength")
		Ineq, Limit, msg = p.ParseLength(cond, dummy)
		fmt.Println("Length: ", Ineq, ',', Limit)
		if Limit == -1 {
			return primitive.NilObjectID, msg
		}
	} else if Name == "Contains" {
		fmt.Println("Going into ParseContains")
		Matcher, tag, msg = p.ParseContains(cond)
		fmt.Println("Contains: ", Matcher)
		if tag == -1 {
			return primitive.NilObjectID, msg
		}
	} else if Name == "MAX" || Name == "MIN" {
		fmt.Println("Going into ParseMinMax")
		Obj, Ineq, Limit, msg = p.ParseMinMax(ctx, cond, dummy)
		fmt.Println(Name, Ineq, ',', Limit)
		if Limit == -1 {
			return primitive.NilObjectID, msg
		}
	} else if Name == "AND" || Name == "OR" {
		fmt.Println("Going into ParseAndOr")
		Obj, tag, msg = p.ParseAndOr(ctx, cond)
		if tag == -1 {
			return primitive.NilObjectID, msg
		}
	} else if Name == "NOT" {
		fmt.Println("Going into ParseNot")
		Obj, tag, msg = p.ParseNot(ctx, cond)
		if tag == -1 {
			return primitive.NilObjectID, msg
		}
	} else {
		msg += fmt.Sprintf("InValid Instruction: %v", Name)
		return primitive.NilObjectID, msg
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
			msg += "Error from CreateRule: " + err.Error()
			return primitive.NilObjectID, msg
		}
	}	
	return ID, "Success"
}

func (p *Parser) ParseRule(ctx context.Context, text string, rule_owner primitive.ObjectID) (Rule, string, error) {
	var rule,empty_rule Rule
	notifyPattern := regexp.MustCompile(`NOTIFY\s+(.*?)\s+WHEN\s+(.*)`)
	matches := notifyPattern.FindStringSubmatch(text)
	if len(matches) != 3 {
		return rule,  "Not of form: NOTIFY \"your_classification\" WHEN \"your_condition\"" , errors.New("not of form: NOTIFY \"your_classification\" WHEN \"your_condition\"")
	}

	ID,msg := p.ParseCondition(ctx, matches[2], 0, matches[1], rule_owner)

	err := p.service.repo.collection.FindOne(ctx, bson.D{{Key: "_id", Value: ID}}).Decode(&rule)

	if ID == primitive.NilObjectID {
		return empty_rule, msg, err
	}

	return rule, "Error in FindOne after Parsing.", err
}
