package rule

import (
	"context"
	//"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
    collection *mongo.Collection
}


func NewRepository(db *mongo.Database, collectionName string) *Repository {
    return &Repository{
        collection: db.Collection(collectionName),
    }
}
// Name, Matcher, Ineq, Notify, When, Obj, Limit, ID, Created)
func (r *Repository) FindDocWithWhenNotify(ctx context.Context, rule *Rule) (*Rule, error) {
    filter := bson.M{
        "notify": rule.Notify,
        "when" : rule.When,
    }
    
    var temp Rule
    err := r.collection.FindOne(context.TODO(),filter).Decode(&temp)
    
    return &temp, err
}

func (r *Repository) Create(ctx context.Context, rule *Rule) error {
    _, err := r.collection.InsertOne(ctx, rule)
    return err
}

func (r *Repository) Get(ctx context.Context, id primitive.ObjectID) (*Rule, error) {
    var rule Rule
    err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&rule)
    return &rule, err
}

func (r *Repository) Update(ctx context.Context, id primitive.ObjectID, rule *Rule) error {
    _, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": rule})
    return err
}

func (r *Repository) Delete(ctx context.Context, id primitive.ObjectID) error {
    _, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
    return err
}

func (r *Repository) PushUserID(ctx context.Context, id primitive.ObjectID, user_id primitive.ObjectID) error {
    _,err := r.collection.UpdateOne(
        ctx,
        bson.M{"_id": id},
        bson.M{"$push": bson.M{"owners": user_id}},
    )
    return err
}
