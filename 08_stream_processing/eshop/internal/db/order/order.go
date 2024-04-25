package order

import (
	"context"
	"eshop/internal/db"
	requestorder "eshop/internal/request/order"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
)

type Status int32

const (
	Status_New Status = iota
	Status_Cancel
	Status_ReservedItems
	Status_ReserveDeliverSlot
	Status_Paid
)

type DB struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func New(db db.Client, dbName, collectName string) *DB {
	d := DB{Client: db.GetClient(), Collection: db.GetClient().Database(dbName).Collection(collectName)}
	return &d
}

func (d *DB) AddNewOrder(ctx context.Context, items []requestorder.Item) (string, *OrderDB, error) {
	dbItems := []Item{}
	for _, item := range items {
		dbItems = append(dbItems, Item{Title: item.Title, Count: item.Count})
	}

	orderDB := OrderDB{
		Status: Status_New,
		Items:  dbItems,
	}
	one, err := d.Collection.InsertOne(ctx, orderDB)
	if err != nil {
		return "", nil, err
	}
	return one.InsertedID.(primitive.ObjectID).Hex(), &orderDB, nil
}

func (d *DB) GetOrder(ctx context.Context, orderId string) (*OrderDB, error) {
	hex, err := primitive.ObjectIDFromHex(orderId)
	if err != nil {
		return nil, err
	}
	one := d.Collection.FindOne(ctx, bson.D{{"_id", hex}})
	if one.Err() != nil {
		return nil, one.Err()
	}
	order := OrderDB{}
	err = one.Decode(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (d *DB) ChangeStatus(ctx context.Context, orderId string, status Status) error {
	hex, err := primitive.ObjectIDFromHex(orderId)
	if err != nil {
		return err
	}
	update := d.Collection.FindOneAndUpdate(ctx, bson.D{{"_id", hex}}, bson.D{{"$set", bson.D{{"status", status}}}})
	if update.Err() != nil {
		return update.Err()
	}

	statusStr := "undefined"
	switch status {
	case Status_ReservedItems:
		statusStr = "ReservedItems"
	case Status_ReserveDeliverSlot:
		statusStr = "ReserveDeliverSlot"
	case Status_Paid:
		statusStr = "Paid"
	case Status_Cancel:
		statusStr = "Cancel"
	}

	if err := d.AddLog(ctx, orderId, fmt.Sprintf("change status =>: %s", statusStr)); err != nil {
		slog.Error("Save log 'reserving success'", "err", err)
	}
	return nil
}

func (d *DB) SetSlot(ctx context.Context, orderId string, slot int32) error {
	hex, err := primitive.ObjectIDFromHex(orderId)
	if err != nil {
		return err
	}
	update := d.Collection.FindOneAndUpdate(ctx, bson.D{{"_id", hex}}, bson.D{{"$set", bson.D{{"slot", slot}}}})
	if update.Err() != nil {
		return update.Err()
	}
	return nil
}

func (d *DB) AddLog(ctx context.Context, orderId string, log string) error {
	hex, err := primitive.ObjectIDFromHex(orderId)
	if err != nil {
		return err
	}
	update := d.Collection.FindOneAndUpdate(ctx, bson.D{{"_id", hex}}, bson.D{{"$push", bson.D{{"log", log}}}})
	if update.Err() != nil {
		return update.Err()
	}
	return nil
}

type OrderDB struct {
	Id     string   `bson:"_id,omitempty"`
	Status Status   `bson:"status"`
	Items  []Item   `bson:"items"`
	Slot   int32    `bson:"slot,omitempty"`
	Log    []string `bson:"log,omitempty"`
}
type Item struct {
	Title string `bson:"title" json:"title"`
	Count int32  `bson:"count" json:"count"`
}
