package global

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var MongoClient *mongo.Database

func InitMongo() {
	clientOptions := options.Client().ApplyURI(Config.Mongo.Address)

	clientOptions.SetMaxPoolSize(uint64(Config.Mongo.MaxPoolSize))

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Mongo is Collection!!!")
	MongoClient = client.Database("wooplus")
}
