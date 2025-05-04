package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alok-mishra143/go-todo/config"
	"github.com/alok-mishra143/go-todo/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main(){

	fmt.Println("Starting server...")

	if err := godotenv.Load(); err != nil {
    log.Println("No .env file found, using system environment variables")
}



		// Connect to MongoDB
		client := config.ConnectDB()
		defer config.CloseDB(client)

	// Fiber instance
	app := fiber.New()

	front_url:=os.Getenv("FRONT_URL")
	if front_url==""{
		front_url="http://localhost:5173"
		fmt.Println("No front url specified, using default url http://localhost:3000")
	}
	
	app.Use(cors.New(cors.Config{
		AllowOrigins: front_url,
		AllowMethods: "GET,POST,PUT,DELETE,PATCH",
	}))

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