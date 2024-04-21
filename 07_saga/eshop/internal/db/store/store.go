package store

import (
	"eshop/internal/db"
	"go.mongodb.org/mongo-driver/mongo"
)

type DB struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func New(db db.Client) *DB {
	d := DB{Client: db.GetClient(), Collection: db.GetClient().Database("Store").Collection("Store")}
	return &d
}

type ItemDB struct {
	Id    string `bson:"_id,omitempty"`
	Title string `bson:"title"`
	Count int32  `bson:"count"`
}
