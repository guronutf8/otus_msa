package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type DB struct {
	Client      *mongo.Client
	dbName      string
	collectName string
}

func New(uri, dbName, collectName string) *DB {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), 2*time.Second)
	defer cancelFunc()
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("%s", uri)).SetRetryReads(true).SetRetryWrites(false)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &DB{Client: client, dbName: dbName, collectName: collectName}
}
func (d *DB) Check() (bool, error) {
	names, err := d.Client.Database(d.dbName).ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		return false, err
	}
	if len(names) != 1 {
		return false, nil
	}
	return true, nil
}

func (d *DB) Init(docs []interface{}) error {
	col := d.Client.Database(d.dbName).Collection(d.collectName)

	_, err := col.InsertMany(context.TODO(), docs)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) GetClient() *mongo.Client {
	return d.Client
}

//func (d *DB) Post(ctx context.Context, user entity.User) (string, error) {
//	collection := d.Client.Database("Users").Collection("Users")
//	one, err := collection.InsertOne(ctx, user)
//	if err != nil {
//		return "", err
//	}
//	if oid, ok := one.InsertedID.(primitive.ObjectID); ok {
//		return oid.Hex(), err
//	}
//	return "", nil
//}
//
//func (d *DB) List(ctx context.Context) ([]entity.User, error) {
//	collection := d.Client.Database("Users").Collection("Users")
//	filter := bson.D{}
//	cur, err := collection.Find(ctx, filter)
//	if err != nil {
//		return nil, err
//	}
//	var res []entity.User
//
//	for cur.Next(ctx) {
//		var u entity.User
//		err := cur.Decode(&u)
//		if err != nil {
//			log.Fatal(err)
//		}
//		res = append(res, u)
//	}
//
//	return res, nil
//}
//
//func (d *DB) Get(ctx context.Context, id string) (*entity.User, error) {
//	collection := d.Client.Database("Users").Collection("Users")
//	objectId, err := primitive.ObjectIDFromHex(id)
//	if err != nil {
//		return nil, err
//	}
//
//	result := collection.FindOne(context.Background(), bson.M{"_id": objectId})
//	user := entity.User{}
//	err = result.Decode(&user)
//	if err != nil {
//		return nil, err
//	}
//
//	return &user, err
//
//	//return nil, nil
//}
