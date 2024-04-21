package delivery

import (
	"eshop/internal/db"
	"go.mongodb.org/mongo-driver/mongo"
)

type DB struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func New(db db.Client) *DB {
	d := DB{Client: db.GetClient(), Collection: db.GetClient().Database("Delivery").Collection("Delivery")}
	return &d
}

type SlotDB struct {
	Id   string `bson:"_id,omitempty"`
	Slot int32  `bson:"slot"`
}
