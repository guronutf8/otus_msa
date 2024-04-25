package billing

import (
	"context"
	"errors"
	"eshop/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
)

type DB struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

type Row struct {
	Id      primitive.ObjectID `bson:"_id,omitempty"`
	User    string             `bson:"user"`
	Balance int32              `bson:"balance"`
}

func New(db db.Client, dbName, collectName string) *DB {
	d := DB{Client: db.GetClient(), Collection: db.GetClient().Database(dbName).Collection(collectName)}
	return &d
}

func (d *DB) CreateUser(ctx context.Context, user string) bool {
	one := d.Collection.FindOne(ctx, bson.D{{"user", user}})
	if one.Err() == nil {
		slog.Warn("User is exist")
		return false
	}
	if !errors.As(one.Err(), &mongo.ErrNoDocuments) {
		slog.Error("Check user, for create", "err", one.Err())
		return false
	}

	userRow := Row{User: user, Balance: 0}
	_, err := d.Collection.InsertOne(ctx, userRow)
	if err != nil {
		slog.Error("Insert user", "err", err)
		return false
	}

	return true
}

func (d *DB) DepositCash(ctx context.Context, user string, sum int32) bool {
	one := d.Collection.FindOne(ctx, bson.D{{"user", user}})
	if one.Err() != nil {
		slog.Warn("DepositCash find")
		return false
	}

	row := &Row{}
	if err := one.Decode(&row); err != nil {
		slog.Error("Deposit cash decode", "err", err)
		return false
	}

	userRow := Row{Id: row.Id, User: user, Balance: row.Balance + sum}
	res := d.Collection.FindOneAndReplace(ctx, bson.D{{"_id", row.Id}}, userRow)
	if res.Err() != nil {
		slog.Error("DepositCash update", "err", res.Err())
		return false
	}

	return true
}

func (d *DB) Pay(ctx context.Context, user string, sum int32) bool {
	one := d.Collection.FindOne(ctx, bson.D{{"user", user}})
	if one.Err() != nil {
		slog.Warn("DepositCash find")
		return false
	}

	row := &Row{}
	if err := one.Decode(&row); err != nil {
		slog.Error("Deposit cash decode", "err", err)
		return false
	}

	if row.Balance < sum {
		slog.Info("User pay fail, little money")
		return false
	}

	userRow := Row{Id: row.Id, User: user, Balance: row.Balance - sum}
	res := d.Collection.FindOneAndReplace(ctx, bson.D{{"_id", row.Id}}, userRow)
	if res.Err() != nil {
		slog.Error("DepositCash update", "err", res.Err())
		return false
	}

	return true
}
