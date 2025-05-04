package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/alok-mishra143/go-todo/handlers"
)

func TodoRoutes(app *fiber.App) {
	api := app.Group("/api")
	todo := api.Group("/todos")

	todo.Get("/", handlers.GetTodos)
	todo.Post("/", handlers.CreateTodo)
	todo.Patch("/:id", handlers.UpdateTodo)
	todo.Delete("/:id", handlers.DeleteTodo)
}
