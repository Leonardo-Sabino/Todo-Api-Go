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

func updateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objetcedID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Id"})
	}

	filter := bson.M{"_id": objetcedID}
	update := bson.M{"$set": bson.M{"completed": true}}

	_, err = colection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": "Todo successfully updated"})

}

func deleteTodos(c *fiber.Ctx) error {
	id := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Todo id"})
	}

	filter := bson.M{"_id": objectID}

	result := colection.FindOne(context.Background(), filter)

	if err := result.Decode(&bson.M{}); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	}

	_, err = colection.DeleteOne(context.Background(), filter)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": "Todo successfully deleted"})
}

func main() {
	fmt.Println("Hello world")

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	//db connection
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
	api.Patch("todos/:id", updateTodo)
	api.Delete("todos/:id", deleteTodos)

	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "3000"
	}

	log.Fatal(app.Listen("0.0.0.0:" + PORT))

}
