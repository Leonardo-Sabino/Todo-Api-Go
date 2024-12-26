package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Todo struct {
	ID        int    `json:"id"`
	COMPLETED bool   `json:"completed"`
	BODY      string `json:"body"`
}

var todos = []Todo{}

func createTodo(c *fiber.Ctx) error {
	newTodo := &Todo{}

	if err := c.BodyParser(newTodo); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if newTodo.BODY == "" {
		return fiber.NewError(fiber.StatusBadRequest, "body is required")
	}

	// Generate a new ID and add the todo to the slice
	newTodo.ID = len(todos) + 1
	todos = append(todos, *newTodo)

	// Return the created todo as a response
	return c.Status(fiber.StatusCreated).JSON(newTodo)
}

func updtadeTodo(c *fiber.Ctx) error {
	id := c.Params("id")

	for i, todo := range todos {
		if fmt.Sprint(todo.ID) == id {
			todos[i].COMPLETED = !todos[i].COMPLETED
			return c.Status(200).JSON(todos[i])
		}

	}
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Todo not found"})
}

func deleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")

	for i, todo := range todos {
		if fmt.Sprint(todo.ID) == id {
			todos = append(todos[:i], todos[i+1:]...)
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Todo successfully deleted"})
		}
	}
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Todo not found"})

}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env variables")
	}

	app := fiber.New()

	api := app.Group("/api")

	PORT := os.Getenv("PORT")

	log.Println("Server is running...")

	api.Get("/todos", func(c *fiber.Ctx) error {
		if len(todos) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No todos found"})
		}
		return c.Status(fiber.StatusOK).JSON(todos)
	})

	api.Post("/todos", createTodo)

	//update a todo
	api.Patch("/todos/:id", updtadeTodo)

	//delete a todo
	api.Delete("/todos/:id", deleteTodo)

	log.Fatal(app.Listen(":" + PORT))
}
