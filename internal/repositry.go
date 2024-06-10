package user

import (
    "context"
    //"log"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    //"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
    collection *mongo.Collection
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
    _, err := r.collection.InsertOne(ctx, user)
    return err
}
