package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	fmt.Println("Starting Fiber Todo API...")

	err:=godotenv.Load()

	if err!=nil {
		log.Fatal("Error loading .env file")
	}

	PORT:=os.Getenv("PORT")

	app := fiber.New()
	todos := []Todo{}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"data": todos,
		})
	})

	// Grouping API routes under /api
	api := app.Group("/api")

	// Create a new todo
	api.Post("/todos", func(c *fiber.Ctx) error {
		todo := new(Todo)

		if err := c.BodyParser(todo); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request payload",
			})
		}

		if todo.Body == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Body is required",
			})
		}

		todo.ID = len(todos) + 1
		todos = append(todos, *todo)

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Todo created",
			"todo":    todo,
		})
	})

	// Update a todo (toggle completion)
	api.Patch("/todos/:id", func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID",
			})
		}

		for i, todo := range todos {
			if todo.ID == id {
				todos[i].Completed = !todos[i].Completed

				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"message": "Todo updated",
					"todo":    todos[i],
				})
			}
		}

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Todo not found",
		})
	})

	// Delete a todo
	api.Delete("/todos/:id", func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid ID",
			})
		}

		for i, todo := range todos {
			if todo.ID == id {
				todos = append(todos[:i], todos[i+1:]...)

				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"message": "Todo deleted",
					"todo":    todo,
				})
			}
		}

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Todo not found",
		})
	})

	// Start the server
	log.Fatal(app.Listen(":"+PORT))
}
