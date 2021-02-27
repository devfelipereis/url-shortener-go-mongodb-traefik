package database

import (
	"context"
	"devfelipereis/urlShortener/env"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() *mongo.Client {
	appEnv := env.Get()

	credentials := options.Credential{
		Username: appEnv.MongoUsername,
		Password: appEnv.MongoPassword,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(context.TODO(), options.Client().SetAuth(credentials).ApplyURI(appEnv.MongoHost))

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
