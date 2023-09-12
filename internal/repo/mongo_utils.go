package repo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Db               *mongo.Database
	MongoClient      *mongo.Client
	ClientRepo       *ClientRepoImpl
	DistributionRepo *DistributionRepoImpl
)

const (
	NUM_CONNECTION_RETRIES     = 3
	distributionCollectionName = "distributions"
	clientCollectionName       = "client"
	messageCollectionName      = "message"
)

func fillCollectionVariables(client *mongo.Client, dbName string) {
	Db = client.Database(dbName)
	DistributionRepo = &DistributionRepoImpl{Db.Collection(distributionCollectionName)}
	ClientRepo = &ClientRepoImpl{Db.Collection(clientCollectionName)}
	messageCollection = Db.Collection(messageCollectionName)
}

func ConnectToMongo(ctx context.Context, uri string, dbName string) {
	for i := 0; i < NUM_CONNECTION_RETRIES; i++ {
		var err error
		MongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			log.Println(err)
			continue
		}
		ctxPing, cancel := context.WithTimeout(ctx, 20*time.Second)
		err = MongoClient.Ping(ctxPing, nil)
		cancel()
		if err != nil {
			log.Println(err)
			continue
		}
		if err != nil {
			log.Println(err)
			continue
		}
		fillCollectionVariables(MongoClient, dbName)
		return
	}
	log.Panic("Connection to mongoDb was not set")
}

func DropDb() {
	DistributionRepo.collection.Drop(context.TODO())
	ClientRepo.collection.Drop(context.TODO())
	messageCollection.Drop(context.TODO())
}
