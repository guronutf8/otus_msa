package notify

import (
	"context"
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
	Id     primitive.ObjectID `bson:"_id,omitempty"`
	User   string             `bson:"user"`
	Result bool               `bson:"result"`
}

func New(db db.Client, dbName, collectName string) *DB {
	d := DB{Client: db.GetClient(), Collection: db.GetClient().Database(dbName).Collection(collectName)}
	return &d
}

func (d *DB) SaveNotify(ctx context.Context, user string, result bool) bool {

	userRow := Row{User: user, Result: result}
	_, err := d.Collection.InsertOne(ctx, userRow)
	if err != nil {
		slog.Error("Insert notify", "err", err)
		return false
	}

	return true
}

func (d *DB) GetLog(ctx context.Context) []Row {

	cursor, err := d.Collection.Find(ctx, bson.D{})
	if err != nil {
		slog.Error("GetLog notify", "err", err)
		return nil
	}

	rows := []Row{}
	err = cursor.All(ctx, &rows)
	if err != nil {
		slog.Error("GetLog cursor all", "err", err)
		return nil
	}

	return rows
}
