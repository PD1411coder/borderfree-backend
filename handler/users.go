package handler

import (
	"borderfree/config"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllUsers(c *fiber.Ctx) error {
	var client = config.ConnectDB()
	var usersCollection *mongo.Collection = config.GetCollection(client, "users")

	cursor, err := usersCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		// client.Disconnect(context.TODO())
		panic(err)
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		// client.Disconnect(context.TODO())
		panic(err)
	}

	// client.Disconnect(context.TODO())
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"payload": results,
	})
}

// err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(payload.Password))
// 	authenticated := true
// 	if err != nil {
// 		authenticated = false
// 	}

// for cursor.Next(context.Background()) {
// 	result := struct{
// 		_id string
// 		username string
// 		password string
// 	}{}

// 	err := cursor.Decode(&result)
// 	if err != nil {
// 		cursor.Close(context.Background())
// 		// client.Disconnect(context.TODO())
// 		return c.Status(500).JSON(&fiber.Map{
// 			"success": false,
// 			"payload": err,
// 		})
// 	}

// 	fmt.Println(result)

// 	if result.username == payload.Username {
// 		cursor.Close(context.Background())
// 		// client.Disconnect(context.TODO())

// 		return c.Status(500).JSON(&fiber.Map{
// 			"success": false,
// 			"payload": "User already registered.",
// 		})
// 	}
// }

// cursor.Close(context.Background())
