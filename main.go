package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/alok-mishra143/go-todo/config"
	"github.com/alok-mishra143/go-todo/routes"

)

func main(){

	fmt.Println("Starting server...")

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}


		// Connect to MongoDB
		client := config.ConnectDB()
		defer config.CloseDB(client)

	// Fiber instance
	app := fiber.New()

	// Register routes
	routes.TodoRoutes(app)

	// Start server
	port := os.Getenv("PORT")

	if port==""{
		port = "4000"
		fmt.Println("No port specified, using default port 4000")
	}
	log.Fatal(app.Listen(":" + port))
}