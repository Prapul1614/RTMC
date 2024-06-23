package rule

import (
	"context"
	"unicode/utf8"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Prapul1614/RTMC/internal/user"
)

type Service struct {
    repo *Repository
	//userRepo *mongo.Collection
    userRepo user.Repository
    mockrepo *MockRepository
}

func NewService(repo *Repository, userRepo user.Repository) *Service {
    return &Service{repo: repo, userRepo : userRepo}
}

func (s* Service) FindDoc(ctx context.Context, rule *Rule, id primitive.ObjectID) (primitive.ObjectID, error) {
    temp, err := s.repo.FindDocWithWhenNotify(ctx, rule) 
    
    if err == nil && id != primitive.NilObjectID{
        for _,v := range temp.Owners {
            if v == id { return temp.ID, err }
        } 

        s.repo.PushUserID(ctx, temp.ID, id)
        s.userRepo.UpdateUser(ctx, id, temp.ID)
    } 
    return temp.ID, err
}

func (s *Service) CreateRule(ctx context.Context, rule *Rule, id primitive.ObjectID) (primitive.ObjectID,error) {
    rule.ID = primitive.NewObjectID()
    s.userRepo.UpdateUser(ctx, id, rule.ID)
    return rule.ID, s.repo.Create(ctx, rule)
}

func (s *Service) GetRule(ctx context.Context, id primitive.ObjectID) ([]Rule, error) {
    // Create the aggregation pipeline
    /*pipeline := mongo.Pipeline{
        bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: id}}}},
        bson.D{{Key: "$unwind", Value: "$rules"}},
        bson.D{{
            Key: "$lookup", Value: bson.D{
                {Key: "from", Value: "rules"},
                {Key: "localField", Value: "rules"},
                {Key: "foreignField", Value: "_id"},
                {Key: "as", Value: "rule_docs"},
            },
        }},
        bson.D{{Key: "$unwind", Value: "$rule_docs"}},
        bson.D{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: "$rule_docs"}}}},
    }

    // Execute the aggregation
    cursor, err := s.userRepo.Collection().Aggregate(ctx, pipeline)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)*/

	cursor, err := s.repo.FindRulesOfUser(ctx, id)
	if err != nil {
		return nil,err
	}
	defer cursor.Close(ctx)

    // Decode the results into a slice of Rule
    var rules []Rule
    if err = cursor.All(ctx, &rules); err != nil {
        return nil, err
    }

    return rules,nil
}

func (s *Service) UpdateRule(ctx context.Context, id primitive.ObjectID, rule *Rule) error {
    return s.repo.Update(ctx, id, rule)
}

func (s *Service) DeleteRule(ctx context.Context, id primitive.ObjectID) error {
    return s.repo.Delete(ctx, id)
}

func GetLength(text string) int {
    return utf8.RuneCountInString(text)
}

const CharSetSize = 256
func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}
func preprocess(pattern string, size int) [CharSetSize]int {
    var badChar [CharSetSize]int
    for i := 0; i < CharSetSize; i++ {
        badChar[i] = -1
    }
    for i := 0; i < size; i++ {
        badChar[int(pattern[i])] = i
    }

    return badChar
}
func GetCount(txt string, pat string, contains bool) int {
    // Boyer-Moore Algorithm
    m := len(pat)
    n := len(txt)
    count := 0
    badChar := preprocess(pat, m)

    s := 0 
    for s <= (n - m) {
        j := m - 1
        for j >= 0 && pat[j] == txt[s+j] {
            j--
        }
        if j < 0 {
            if contains {return 1}
            count++
            if s+m < n {
                s += max(1, m-badChar[txt[s+m]])
            } else { break }

        } else {
            s += max(1, j-badChar[txt[s+j]])
        }
    }
    return count
}

func GetContains(text string, word string) bool {
    count := GetCount(text, word, true)
    return count == 1
}
func ImplementOp(op string, ans, limit int) bool {
    if ans == -1 { return false}
    if(op == "<=") { return ans <= limit
    } else if op == "<" { return ans < limit
    } else if op == "=" { return ans == limit
    } else if op == "!=" { return ans != limit
    } else if op == ">" { return ans > limit
    } else if op == ">=" { return ans >= limit }
    return false
}

func (s *Service) ImplementMinMax(ctx context.Context,text string, rule *Rule) int {
    switch rule.Name {
    case "Length":
        return GetLength(text)
    case "Count":
        return GetCount(text, rule.Matcher, false)
    case "MIN", "MAX":
        ans := -1 
        if rule.Name == "MIN" {ans = 2147483647}
        var Nrule *Rule
        var err error
        for _,v := range rule.Obj {
            if(s.mockrepo != nil) {Nrule, err = s.mockrepo.Get(ctx, v) 
            }else {Nrule, err = s.repo.Get(ctx, v)}
            
            if err != nil{
                println("Can't get rule doc with rule.MinMax.Obj", rule.ID.Hex(), v.Hex())
                return -1
            }
            if rule.Name == "MIN" { 
                ans = min(ans, s.ImplementMinMax(ctx, text, Nrule))
            } else {
                ans = max(ans, s.ImplementMinMax(ctx, text, Nrule))
            }
        }
        return ans
    default:
        println("UNKNOWN RULE IN MINMAX")
        return -1
    }
}
func (s *Service) ImplementAndOr(ctx context.Context, text string, rule *Rule) bool {
    var ans = true
    if rule.Name == "OR" {ans = false}
    var Nrule *Rule
    var err error
    for _,v := range rule.Obj {
        if(s.mockrepo != nil) {Nrule, err = s.mockrepo.Get(ctx, v) 
        }else {Nrule, err = s.repo.Get(ctx, v)}

        if err != nil{
            println("Can't get rule doc with ruleId",v.Hex(),"in", rule.ID.Hex() )
            return false
        }
        if rule.Name == "AND" { 
            ans = ans && s.ImplementRule(ctx, text, Nrule)
        } else {
            ans = ans || s.ImplementRule(ctx, text, Nrule)
        }
    }
    return ans
}
func (s *Service) ImplementNot(ctx context.Context, text string, rule *Rule) bool {
    Nrule, err := s.repo.Get(ctx, rule.Obj[0])
    if err != nil{
        println("Can't get rule doc with ruleId",rule.Obj[0].Hex(),"in", rule.ID.Hex() )
        return false
    }
    return !s.ImplementRule(ctx, text, Nrule)
}

func (s *Service) ImplementRule(ctx context.Context, text string, rule *Rule) bool {
    switch rule.Name {
    case "Length":
        ans := GetLength(text)
        return ImplementOp(rule.Ineq, ans, rule.Limit)
    case "Count":
        ans := GetCount(text, rule.Matcher, false)
        return ImplementOp(rule.Ineq, ans, rule.Limit)
    case "Contains":
        return GetContains(text, rule.Matcher)
    case "MIN", "MAX":
        ans := s.ImplementMinMax(ctx, text, rule)
        return ImplementOp(rule.Ineq, ans, rule.Limit)
    case "AND", "OR":
        return s.ImplementAndOr(ctx, text, rule)
    case "NOT":
        return s.ImplementNot(ctx, text, rule)
    default:
        println("UNKNOWN RULE")
    }
    return true
}

