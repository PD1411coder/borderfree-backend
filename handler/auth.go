package handler

import (
	"borderfree/config"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx) error {
	payload := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"payload": err.Error(),
		})
	}

	var client = config.ConnectDB()
	var usersCollection *mongo.Collection = config.GetCollection(client, "users")

	result := struct {
		Id       primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
		Username string             `json:"username"`
		Password string             `json:"password"`
	}{}
	err := usersCollection.FindOne(context.TODO(), bson.M{"username": payload.Username}).Decode(&result)
	// client.Disconnect(context.TODO())
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"success": true,
			"payload": "User is not registered!",
		})
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(payload.Password))
		if err != nil {
			return c.Status(fiber.StatusOK).JSON(&fiber.Map{
				"success": false,
				"payload": "Invalid credentials!",
			})
		} else {
			claims := jwt.MapClaims{
				"username": result.Username,
				"userid":   result.Id,
				"exp":      time.Now().Add(time.Hour * 72).Unix(),
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			t, err := token.SignedString([]byte("ANiceLittleSecretOfMine987654321"))
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
					"success": false,
					"payload": err.Error(),
				})
			}

			return c.Status(fiber.StatusOK).JSON(&fiber.Map{
				"success": true,
				"payload": "Logged in.",
				"token":   t,
			})
		}
	}
}

func Register(c *fiber.Ctx) error {
	payload := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"payload": err.Error(),
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"payload": err.Error(),
		})
	}

	// Checking if this user is already registered.
	var client = config.ConnectDB()
	var usersCollection *mongo.Collection = config.GetCollection(client, "users")

	var result bson.M
	err = usersCollection.FindOne(context.TODO(), bson.M{"username": payload.Username}).Decode(&result)
	if err != nil {
		_, err = usersCollection.InsertOne(context.TODO(), bson.M{
			"username": payload.Username,
			"password": hashedPassword,
		})
		if err != nil {
			// client.Disconnect(context.TODO())
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"success": false,
				"payload": err.Error(),
			})
		}

		// client.Disconnect(context.TODO())
		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"success": true,
			"payload": "User registered successfully.",
		})
	} else {
		// client.Disconnect(context.TODO())
		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"success": false,
			"payload": "User is already registered!",
		})
	}
}
