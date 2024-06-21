package rule

import (
	"context"
	"time"
	//"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
    collection *mongo.Collection
}

func (r *Repository) Collection() *mongo.Collection {
    return r.collection
}

func NewRepository(db *mongo.Database, collectionName string) *Repository {
    return &Repository{
        collection: db.Collection(collectionName),
    }
}

func (r *Repository)CreateIndexOwners() error{
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "owners", Value: 1}}, // Index on the BSON field name "owners"
	}
	_, err := r.collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}
    return nil
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

func (r *Repository) Aggregate(ctx context.Context, pipeline mongo.Pipeline) (*mongo.Cursor, error){
    curser,err := r.collection.Aggregate(ctx, pipeline)
    return curser, err
}

func (r *Repository) FindRulesOfUser(ctx context.Context, id primitive.ObjectID) (*mongo.Cursor, error){
    filter := bson.M{"owners": id}
    cursor, err := r.collection.Find(ctx, filter)
    return cursor, err
}
