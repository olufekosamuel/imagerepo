package helpers

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	err      = godotenv.Load()
	mongoUri = os.Getenv("MONGO_DB")
	DB, _    = InitDB()
)

func InitDB() (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		mongoUri,
	))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	db := client.Database("imagerepo")
	return db, nil
}

func GetDBCollection(name string) (*mongo.Collection, error) {

	client := DB

	collection := client.Collection(name)
	return collection, nil
}
