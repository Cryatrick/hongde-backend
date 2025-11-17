package database


import (
	"log"
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"hongde_backend/internal/config"
)

var DbMongo *mongo.Database

func OpenMongo() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	uri := "mongodb://" + config.DB_MONGO_HOSTNAME
	if config.DB_MONGO_PASSWORD != "" {
		uri = "mongodb://"+config.DB_MONGO_USERNAME+":"+config.DB_MONGO_PASSWORD+"@"+config.DB_MONGO_HOSTNAME+"/?authSource=admin"
	}

	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(opts)

	if err != nil {
		log.Fatalf("Failed to connect to DB MONGO %v", err)
	}
	
	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatalf("Failed to ping DB MONGO %v", err)
	}

	DbMongo = client.Database(config.DB_MONGO_DBNAME)
}