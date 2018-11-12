package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/CastellanR/UserFeedback-Microservice/tools/env"
	"github.com/mongodb/mongo-go-driver/core/topology"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
)

var database *mongo.Database

// Get the mongo database
func Get() (*mongo.Database, error) {
	if database == nil {
		client, err := mongo.NewClientWithOptions(
			env.Get().MongoURL,
			clientopt.ServerSelectionTimeout(time.Second))
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		err = client.Connect(context.TODO())
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		database = client.Database("feedback")
		fmt.Println("bd")
	}
	return database, nil
}

// CheckError función a llamar cuando se produce un error de db
func CheckError(err interface{}) {
	if err == topology.ErrServerSelectionTimeout {
		database = nil
	}
}

// IsUniqueKeyError retorna true si el error es de indice único
func IsUniqueKeyError(err error) bool {
	if wErr, ok := err.(mongo.WriteErrors); ok {
		for i := 0; i < len(wErr); i++ {
			if wErr[i].Code == 11000 {
				return true
			}
		}
	}
	return false
}
