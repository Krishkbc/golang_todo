package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/thedevsaddam/renderer"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
	checkerr(err)

	err = client.Ping(ctx, readpref.Primary())
	checkerr(err)

	db = client.Database(dbname)

}

func main() {
	fmt.Println("helloo")
	server := &http.Server{
		Addr:         ":9000",
		Handler:      chi.NewRouter(),
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	fmt.Println("server starting at port 9000")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
