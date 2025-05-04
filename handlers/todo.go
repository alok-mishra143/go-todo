package handlers

import (
	"time"

	"github.com/alok-mishra143/go-todo/config"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Completed bool               `json:"completed"`
	Body      string             `json:"body"`
	CreatedAt primitive.DateTime `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}



// ? Get all todos

func GetTodos(c *fiber.Ctx) error {
	db := config.GetCollection()
	var todo []Todo

	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor,err:=db.Find(c.Context(),bson.M{},opts)

	if err!=nil{
		return c.Status(500).JSON(fiber.Map{"error": "Error fetching todos"})
	}

	for cursor.Next(c.Context()){
		var t Todo
		if err := cursor.Decode(&t); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error decoding todo"})
		}
		todo = append(todo, t)
	}

	if err := cursor.Close(c.Context()); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error closing cursor"})
	}
	return c.JSON(fiber.Map{"data": todo})
}

// ? Create a new todo

func CreateTodo(c *fiber.Ctx) error {
	var todo =new(Todo)

	if err := c.BodyParser(&todo); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Error parsing todo"})
	}

	if todo.Body==""{
		return c.Status(400).JSON(fiber.Map{"error": "Body is required"})
	}

	db := config.GetCollection()
	todo.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	Result,err:=db.InsertOne(c.Context(),todo)
	if err!=nil{
		return c.Status(500).JSON(fiber.Map{"error": "Error inserting todo"})
	}

	todo.ID=Result.InsertedID.(primitive.ObjectID)


	return c.Status(201).JSON(fiber.Map{"meassage": "Todo created successfully", "todo": todo})
}

func UpdateTodo(c *fiber.Ctx) error {
	db := config.GetCollection()

	// Get and validate the ObjectID
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid todo ID",
		})
	}

	// Find the existing todo
	var todo Todo
	err = db.FindOne(c.Context(), bson.M{"_id": objectID}).Decode(&todo)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Todo not found",
		})
	}

	// Toggle the completed field
	newStatus := !todo.Completed
	update := bson.M{"$set": bson.M{"completed": newStatus}}

	// Apply the update
	_, err = db.UpdateOne(c.Context(), bson.M{"_id": objectID}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update todo",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Todo updated successfully",
		"completed": newStatus,
	})
}

func DeleteTodo(c *fiber.Ctx) error {
	db := config.GetCollection()

	// Get and validate the ObjectID
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid todo ID",
		})
	}

	// Attempt to delete the document
	result, err := db.DeleteOne(c.Context(), bson.M{"_id": objectID})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete todo",
		})
	}

	// Check if a document was actually deleted
	if result.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Todo not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Todo deleted successfully",
	})
}
