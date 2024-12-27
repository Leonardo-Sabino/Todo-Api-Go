package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	COMPLETED bool               `json:"completed"`
	BODY      string             `json:"body"`
}

var colection *mongo.Collection

func getTodos(c *fiber.Ctx) error {
	var todos []Todo

	cursor, err := colection.Find(context.Background(), bson.M{})

	if err != nil {
		return err
	}

	defer cursor.Close(context.Background()) // defer is used to postpone the execution of somenthing

	for cursor.Next(context.Background()) {
		var todo Todo

		if err := cursor.Decode(&todo); err != nil {
			return err
		}

		todos = append(todos, todo)
	}
	return c.JSON(todos)

}

func createTodo(c *fiber.Ctx) error {
	todo := new(Todo)

	if err := c.BodyParser(todo); err != nil {
		return err
	}

	if todo.BODY == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Todo's body cannot be empty"})
	}

	insertResult, err := colection.InsertOne(context.Background(), todo)

	if err != nil {
		return err
	}

	todo.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(fiber.StatusCreated).JSON(todo)
}

// func updateTodo (c * fiber.Ctx) error {

// }
// func deleteTodos (c * fiber.Ctx) error {

// }

func main() {
	fmt.Println("Hello world")

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")

	clientOptions := options.Client().ApplyURI(MONGODB_URI)

	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MONGODB ATLAS")

	colection = client.Database("golang_db").Collection("todos")

	app := fiber.New()

	api := app.Group("/api")

	api.Get("todos", getTodos)
	api.Post("todos", createTodo)
	// api.Patch("todos", updateTodo)
	// api.Delete("todos", deleteTodos)

	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "3000"
	}

	log.Fatal(app.Listen("0.0.0.0:" + PORT))

}
