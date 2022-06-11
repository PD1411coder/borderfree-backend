package handler

import (
	"borderfree/config"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllProducts(c *fiber.Ctx) error {
	var client = config.DB
	var productsCollection *mongo.Collection = config.GetCollection(client, "products")

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["userid"].(string)

	parsedUserId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		// client.Disconnect(context.TODO())
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"payload": err.Error(),
		})
	}

	cursor, err := productsCollection.Find(context.TODO(), bson.M{"owner": parsedUserId})
	if err != nil {
		// client.Disconnect(context.TODO())
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"payload": err.Error(),
		})
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		// client.Disconnect(context.TODO())
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"payload": err.Error(),
		})
	}

	// client.Disconnect(context.TODO())
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"payload": results,
	})
}

func DeleteProduct(c *fiber.Ctx) error {
	payload := struct {
		ProductId string `json:"productId"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"payload": err.Error(),
		})
	}

	var client = config.DB
	var productsCollection *mongo.Collection = config.GetCollection(client, "products")

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["userid"].(string)

	parsedUserId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		// client.Disconnect(context.TODO())
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"payload": err.Error(),
		})
	}

	parsedProductId, err := primitive.ObjectIDFromHex(payload.ProductId)
	if err != nil {
		// client.Disconnect(context.TODO())
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"payload": err.Error(),
		})
	}

	targetProduct := struct {
		Id    primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
		Owner primitive.ObjectID `bson:"owner" json:"owner,omitempty"`
	}{}
	err = productsCollection.FindOne(context.TODO(), bson.M{"_id": parsedProductId}).Decode(&targetProduct)
	if err != nil {
		// client.Disconnect(context.TODO())
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"payload": err.Error(),
		})
	}

	if targetProduct.Owner == parsedUserId {
		productsCollection.DeleteOne(context.TODO(), bson.M{"_id": parsedProductId})
		// client.Disconnect(context.TODO())
		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"success": true,
			"payload": "Product with ID " + payload.ProductId + " deleted successfully.",
		})
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"success": false,
			"payload": "You do not own this product!",
		})
	}
}

func AddProduct(c *fiber.Ctx) error {
	payload := struct {
		ProductName  string `json:"productName"`
		ProductPrice int    `json:"productPrice"`
		ProductType  string `json:"productType"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"payload": err.Error(),
		})
	}

	var client = config.DB
	var productsCollection *mongo.Collection = config.GetCollection(client, "products")

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["userid"].(string)

	parsedUserId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		// client.Disconnect(context.TODO())
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"payload": err.Error(),
		})
	}

	_, _ = productsCollection.InsertOne(context.TODO(), bson.M{
		"productName":  payload.ProductName,
		"productPrice": payload.ProductPrice,
		"productType":  payload.ProductType,
		"owner":        parsedUserId,
	})

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"success": true,
		"payload": "Product added successfully.",
	})
}

func UpdateProduct(c *fiber.Ctx) error {
	payload := struct {
		ProductId    string `json:"productId"`
		ProductName  string `json:"productName"`
		ProductPrice int    `json:"productPrice"`
		ProductType  string `json:"productType"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"payload": err.Error(),
		})
	}

	var client = config.DB
	var productsCollection *mongo.Collection = config.GetCollection(client, "products")

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := claims["userid"].(string)

	parsedUserId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		// client.Disconnect(context.TODO())
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"payload": err.Error(),
		})
	}

	parsedProductId, err := primitive.ObjectIDFromHex(payload.ProductId)
	if err != nil {
		// client.Disconnect(context.TODO())
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"payload": err.Error(),
		})
	}

	targetProduct := struct {
		Id    primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
		Owner primitive.ObjectID `bson:"owner" json:"owner,omitempty"`
	}{}
	err = productsCollection.FindOne(context.TODO(), bson.M{"_id": parsedProductId}).Decode(&targetProduct)
	if err != nil {
		// client.Disconnect(context.TODO())
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"payload": err.Error(),
		})
	}

	if targetProduct.Owner == parsedUserId {
		_, err = productsCollection.UpdateOne(context.TODO(), bson.M{"_id": parsedProductId}, bson.M{
			"$set": bson.M{
				"productName":  payload.ProductName,
				"productType":  payload.ProductType,
				"productPrice": payload.ProductPrice,
			},
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
			"payload": "Product with ID " + payload.ProductId + " updated successfully.",
		})
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"success": false,
			"payload": "You do not own this product!",
		})
	}
}
