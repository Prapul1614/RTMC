package rule

import (
	"context"
	//"fmt"
	//"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

    "github.com/Prapul1614/RTMC/internal/user"
)

type Service struct {
    repo *Repository
	//userRepo *mongo.Collection
    userRepo user.Repository
}

func NewService(repo *Repository, userRepo user.Repository) *Service {
    return &Service{repo: repo, userRepo : userRepo}
}

func (s* Service) FindDoc(ctx context.Context, rule *Rule, id primitive.ObjectID) (*Rule, error) {
    temp, err := s.repo.FindDoc(ctx, rule) 
    for _,v := range temp.Owners {
        if v == id {
            println("Rule Already Submitted")
            return temp, err
        }
    } 
    
    if err == nil{
        s.repo.PushUserID(ctx, temp.ID, id)
        s.userRepo.UpdateUser(ctx, id, temp.ID)
    } 
    return temp, err
}

func (s *Service) CreateRule(ctx context.Context, rule *Rule, id primitive.ObjectID) error {
    rule.ID = primitive.NewObjectID()
    s.userRepo.UpdateUser(ctx, id, rule.ID)
    return s.repo.Create(ctx, rule)
}

func (s *Service) GetRule(ctx context.Context, id primitive.ObjectID) ([]Rule, error) {
    // Create the aggregation pipeline
    pipeline := mongo.Pipeline{
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
