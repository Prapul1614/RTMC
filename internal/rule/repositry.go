package rule

import (
    "context"

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

func (r *Repository) FindDoc(ctx context.Context, rule *Rule) (*Rule, error) {
    filter := bson.M{
        "name": rule.Name,
        "word": rule.Word,
        "char": rule.Char,
        "limit": rule.Limit,
        "is_their": rule.IsTheir,
        "ineq": rule.Ineq,
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
