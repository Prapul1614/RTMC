package user

import (
	"context"
	//"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
    collection *mongo.Collection
}

// Provide a method to access the collection
func (r *Repository) Collection() *mongo.Collection {
    return r.collection
}

func NewRepository(db *mongo.Database, collectionName string) *Repository {
    return &Repository{
        collection: db.Collection(collectionName),
    }
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (*User, error) {
    var user User
    filter := bson.M{"username": username}
    err := r.collection.FindOne(ctx, filter).Decode(&user)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *Repository) InsertUser(ctx context.Context, user *User) error {
    user.Rules = []primitive.ObjectID{}
    _, err := r.collection.InsertOne(ctx, user)
    return err
}

func (r *Repository) UpdateUser(ctx context.Context, user_id primitive.ObjectID, rule_id primitive.ObjectID) {
    println(rule_id.Hex())
    _,err1 := r.collection.UpdateOne(
        ctx,
        bson.M{"_id": user_id},
        bson.M{"$push": bson.M{"rules": rule_id}},
    )
    if err1 != nil{
        println("Unable to push", err1.Error())
    }
}
