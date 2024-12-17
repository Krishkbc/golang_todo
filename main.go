package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/thedevsaddam/renderer"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var rnd *renderer.Render
var client *mongo.Client
var db *mongo.Database

const (
	dbname         string = "golang-todo"
	collectionName string = "todo"
)

type (
	TodoModel struct {
		ID        string    `json:"id"`
		Title     string    `json:"title"`
		Completed bool      `json:"completed"`
		CreatedAt time.Time `json:"created_at"`
	}

	Todo struct {
		ID        primitive.ObjectID `bson:"id"`
		Title     string             `bson:"title"`
		Completed bool               `bson:"completed"`
		CreatedAt time.Time          `bson:"creation_at"`
	}
)

func checkerr(err error) {
	if err != nil {
		log.Fatal(err)

	}
}

func init() {
	fmt.Println("starting the application")

	rnd = renderer.New()
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	checkerr

}
