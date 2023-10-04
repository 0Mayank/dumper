package configs

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var lock = &sync.Mutex{}
var db *mongo.Client

func ConnectDB() {
	lock.Lock()
	defer lock.Unlock()

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var err error
	db, err = mongo.Connect(ctx, options.Client().ApplyURI(GetConfig().MongoDBUri))
	if err != nil {
		log.Fatal(err)
	}

	// ping the database
	err = db.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
}

func GetDB() *mongo.Client {
	if db == nil {
		ConnectDB()
	}
	return db
}

// getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database(GetConfig().Database).Collection(collectionName)
	return collection
}
