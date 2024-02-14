package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
	"usercrud/internal/entity"
)

type DB struct {
	Client *mongo.Client
}

func New(url, user, password string) *DB {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancelFunc()
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:27017/?authMechanism=SCRAM-SHA-1", user, password, url)) //.SetRetryReads(true).SetRetryWrites(true)
	//clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://root:root@localhost:27017/?authMechanism=SCRAM-SHA-1")) //.SetRetryReads(true).SetRetryWrites(true)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &DB{Client: client}
}

func (d *DB) Post(ctx context.Context, user entity.User) (string, error) {
	collection := d.Client.Database("Users").Collection("Users")
	one, err := collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}
	if oid, ok := one.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), err
	}
	return "", nil
}

func (d *DB) List(ctx context.Context) ([]entity.User, error) {
	collection := d.Client.Database("Users").Collection("Users")
	filter := bson.D{}
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var res []entity.User

	for cur.Next(ctx) {
		var u entity.User
		err := cur.Decode(&u)
		if err != nil {
			log.Fatal(err)
		}
		res = append(res, u)
	}

	return res, nil
}

func (d *DB) Get(ctx context.Context, id string) (*entity.User, error) {
	collection := d.Client.Database("Users").Collection("Users")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result := collection.FindOne(context.Background(), bson.M{"_id": objectId})
	user := entity.User{}
	err = result.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, err

	//return nil, nil
}
