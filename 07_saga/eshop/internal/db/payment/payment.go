package payment

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
