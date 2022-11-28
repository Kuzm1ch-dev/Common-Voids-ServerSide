package controller

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Message struct {
	Text   string
	Length int
}

const connectionString = "mongodb://localhost:27017"
const userDB = "users"

var collection *mongo.Collection

func Save(m *Message) {
	clientOption := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Mongo connect!")

	collection = client.Database(userDB).Collection("messageHistory")
	insertResult, err := collection.InsertOne(context.TODO(), m)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	fmt.Println("Collection OK!")
}
