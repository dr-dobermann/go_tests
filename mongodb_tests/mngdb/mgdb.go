package main

/*
based on docker image created as followed

docker run -d --name=mongo -p 27017-27019:27017-27019 mongo:latest
*/

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// TaxPayer consists data about single tax payer
type TaxPayer struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string `bson:"name"`
	TIN  string `bson:"iin"`
	City string `bson:"city"`
}

func main() {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("test").Collection("taxpayers")

	var tp = TaxPayer {
		Name: "Ruslan Gabitov",
		TIN: "730223303485",
		City: "Almaty",		
	}

	res, err := collection.InsertOne(context.Background(), tp)
	if err != nil {
		log.Fatalf("could not inset a taxpayer: %v", err)
	}
	id := res.InsertedID

	fmt.Printf("Taxpayer successfully added with id %v", id)

	cur, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Fatalf("couldn't get a cursor: %v", err)
	}

	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		err := cur.Decode(&tp)
		if err != nil {
			log.Fatalf("could not decode a cursor record: %v", err)
		}

		fmt.Printf("Taxpayer info:\n  ID: %s\n  Name: %s\n  IIN: %s\n  City: %s\n", tp.ID.Hex(), tp.Name, tp.TIN, tp.City)
	}
	if err := cur.Err(); err != nil {
		log.Fatalf("error while fetching taxpayers info: %v", err)
	}
}