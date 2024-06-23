package rule

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockRepository struct {}

var ruleIDs [10]primitive.ObjectID
var mockRules map[primitive.ObjectID]*Rule

// Initialize the rule IDs and mockRules in the init function
func init() {
    // Generate random ObjectIDs
    for i := 0; i < len(ruleIDs); i++ {
        ruleIDs[i] = primitive.NewObjectID()
    }

    // Initialize the mockRules map
    mockRules = map[primitive.ObjectID]*Rule{
        ruleIDs[0]: {Name: "Length"},
        ruleIDs[1]: {Name: "Count", Matcher: "a"},
        ruleIDs[2]: {Name: "MIN", Obj: []primitive.ObjectID{ruleIDs[0], ruleIDs[1]}},
		ruleIDs[3]: {Name: "Count", Matcher: "b"},
        ruleIDs[4]: {Name: "MAX", Obj: []primitive.ObjectID{ruleIDs[2], ruleIDs[3]}},
        ruleIDs[5]: {Name: "Length", Ineq: ">", Limit: 10},
        ruleIDs[6]: {Name: "Contains", Matcher: "ha"},
        ruleIDs[7]: {Name: "OR", Obj: []primitive.ObjectID{ruleIDs[5], ruleIDs[6]}},
		ruleIDs[8]: {Name: "MAX", Obj: []primitive.ObjectID{ruleIDs[1], ruleIDs[3]}, Ineq: "<", Limit: 5},
        ruleIDs[9]: {Name: "AND", Obj: []primitive.ObjectID{ruleIDs[7], ruleIDs[8]}},
    }
}

func (m *MockRepository) Get(ctx context.Context, id primitive.ObjectID) (*Rule, error) {
    rule, exists := mockRules[id]
    if !exists {
        return nil, errors.New("rule not found")
    }
    return rule, nil
}
