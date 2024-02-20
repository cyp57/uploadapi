package mongodb

import (
	"context"
	"log"

	"time"

	"github.com/cyp57/uploadapi/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Database *mongo.Database

func dbConnection(c *mongo.Database) {
	Database = c
}

func MongoDbConnect(dbConfig config.IDbConfig) {

	credential := options.Credential{
		Username: dbConfig.DbUser(),
		Password: dbConfig.DbPassword(),
	}

	clientOptions := options.Client().ApplyURI(dbConfig.DbHost()).SetAuth(credential)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalln(err)
	}

	//Cancel context to avoid memory leak
	defer cancel()
	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("DB Connected!")
	}

	db := client.Database(dbConfig.DbName())
	dbConnection(db)

}
