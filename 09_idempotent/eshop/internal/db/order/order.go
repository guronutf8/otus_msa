package order

import (
	"context"
	"errors"
	"eshop/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Status int32

const (
	Status_New Status = iota
	Status_Paid
)

const collBasket = "Basket"

type DB struct {
	Client     *mongo.Client
	CollBasket *mongo.Collection
}

func New(db db.Client, dbName, collectName string) *DB {
	d := DB{Client: db.GetClient(), CollBasket: db.GetClient().Database(dbName).Collection(collBasket)}
	return &d
}

// GetCurrentBasket ищет в бд корзину со статусом new, если ее нет то создаст
func (d *DB) GetCurrentBasket(ctx context.Context, user string) (string, error) {
	one := d.CollBasket.FindOne(ctx, BasketDB{Status: Status_New, User: user})
	if one.Err() != nil {
		if errors.Is(one.Err(), mongo.ErrNoDocuments) {
			one, err := d.CollBasket.InsertOne(ctx, BasketDB{Status: Status_New, User: user})
			if err != nil {
				return "", err
			}

			if oid, ok := one.InsertedID.(primitive.ObjectID); ok {
				return oid.Hex(), nil
			}
			//!ok не обработал, но хрен с ним

		}
		return "", one.Err()
	}
	basketDB := BasketDB{}
	if err := one.Decode(&basketDB); err != nil {
		return "", err
	}

	return basketDB.Id, nil

}

func (d *DB) ChangeStatus(ctx context.Context, user, basketId string) (bool, error) {
	basketIdHex, err := primitive.ObjectIDFromHex(basketId)
	if err != nil {
		return false, err
	}

	one := d.CollBasket.FindOne(ctx, bson.D{{"user", user}, {"_id", basketIdHex}})
	if one.Err() != nil {
		return false, err
	}

	basket := BasketDB{}
	err = one.Decode(&basket)
	if one.Err() != nil {
		return false, err
	}

	// уже заказан
	if basket.Status == Status_Paid {
		return false, nil
	}

	update := d.CollBasket.FindOneAndUpdate(ctx, bson.D{{"user", user}, {"status", Status_New}, {"_id", basketIdHex}}, bson.D{{"$set", bson.D{{"status", Status_Paid}}}})
	if update.Err() != nil {
		return false, update.Err()
	}
	return true, nil
}

type BasketDB struct {
	Id     string `bson:"_id,omitempty"`
	Status Status `bson:"status"`
	User   string `bson:"user"`
}
