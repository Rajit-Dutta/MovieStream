package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectDB() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Print("ENV file is missing")
	}

	MongoDb := os.Getenv("MONGO_URI")
	if MongoDb == "" {
		log.Fatal("MONGO_URL is missing")
	}

	clientOptions := options.Client().ApplyURI(MongoDb)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil
	}

	return client
}

var Client *mongo.Client = ConnectDB()

func OpenCollection(collectionName string) *mongo.Collection {
	err := godotenv.Load(".env")
	if err != nil {
		log.Print("ENV file is missing")
	}

	databaseName := os.Getenv("DATABASE_NAME")
	collection := Client.Database(databaseName).Collection(collectionName)

	if collection == nil {
		return nil
	}

	return collection
}
