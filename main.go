package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/thedevsaddam/renderer"
	"go.mongodb.org/mongo-driver/bson"
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
		ID        primitive.ObjectID `json:"id"`
		Title     string             `json:"title"`
		Completed bool               `json:"completed"`
		CreatedAt time.Time          `json:"created_at"`
	}

	Todo struct {
		ID        primitive.ObjectID `bson:"id"`
		Title     string             `bson:"title"`
		Completed bool               `bson:"completed"`
		CreatedAt time.Time          `bson:"creation_at"`
	}

	GetTodoResponse struct {
		Message string `json:"message"`
		Data    []Todo `json:"data"`
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

func homeHandler(w http.ResponseWriter, r *http.Request) {
	filepath := "./README.md"
	err := rnd.FileView(w, http.StatusOK, filepath, "readme.md")
	checkerr(err)
}

func main() {
	fmt.Println("helloo")

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", homeHandler)
	router.Mount("/todo", todoHandler())

	server := &http.Server{
		Addr:         ":9000",
		Handler:      chi.NewRouter(),
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	stopchan := make(chan os.Signal, 1)
	signal.Notify(stopchan, os.Interrupt)

	go func() {
		fmt.Println("server starting at port 9000")
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	sig := <-stopchan
	log.Printf("shutting down server with signal: %v", sig)

	if err := client.Disconnect(context.Background()); err != nil {
		panic(err)
	}
	// create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v\n", err)
	}
	log.Println("Server shutdown gracefully")

}

func todoHandler() http.Handler {
	router := chi.NewRouter()
	router.Group(func(r chi.Router) {
		r.Get("/", getTodos)
		r.Post("/", createTodo)
		r.Put("/{id}", updateTodo)
		r.Delete("/{id}", deleteTodo)
	})
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	var todoListFromDb = []TodoModel{}
	filter := bson.D{}

	cursor, err := db.Collection(collectionName).Find(context.Background(), filter)
	if err != nil {
		log.Printf("failed to fetch todo records from the db: %v\n", err.Error())
		rnd.JSON(w, http.StatusBadRequest, renderer.M{
			"message": "Could not fetch the todo collection",
			"error":   err.Error(),
		})
	}

	todoList := []Todo{}
	if err := cursor.All(context.Background(), &todoListFromDb); err != nil {
		checkerr(err)
	}

	for _, td := range todoListFromDb {
		todoList = append(todoList, Todo{
			ID:        td.ID,
			Title:     td.Title,
			Completed: td.Completed,
			CreatedAt: td.CreatedAt,
		})
	}

	rnd.JSON(w, http.StatusOK, GetTodoResponse{
		Message: "Todo list fetched successfully",
		Data:    todoList,
	})

}
