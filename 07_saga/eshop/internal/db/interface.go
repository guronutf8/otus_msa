package db

import "go.mongodb.org/mongo-driver/mongo"

type Client interface {
	GetClient() *mongo.Client
}
