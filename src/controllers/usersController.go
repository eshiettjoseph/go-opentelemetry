package controllers

import (
	"go-opentelemetry/initializers"
	"go-opentelemetry/models"
	"log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"github.com/gin-gonic/gin"
)

// Get data off req body
// Create a user
// return it
func Status(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "running",
	})
}
func AddUser(c *gin.Context) {

	// struct to store user data
	var body struct {
		ID uint 
		Name string 
		Email *string 
		Age uint8
	}
	c.Bind(&body)

	// create a user 
	user := models.User{ID: body.ID, Name: body.Name, Email: body.Email, Age: body.Age}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(400, gin.H{
			"user": "Bad request",
		})
		return
	}

	// return created user
	c.JSON(200, gin.H{
		"user": user,
	})
}

func GetUsers(c *gin.Context) {
	// Get users and respond with users
	var users []models.User
	span := trace.SpanFromContext(c.Request.Context())
	span.SetAttributes(attribute.String("controller", "users"))
	span.AddEvent("This is a sample event", trace.WithAttributes(attribute.Int("pid", 4328), attribute.String("sampleAttribute", "Test")))
	result := initializers.DB.Find(&users)
	if result.Error != nil{
		log.Fatal("Unable to get users:"+ result.Error.Error())
	}
	c.JSON(200, gin.H{
		"user": users,
	})
}

func GetUser(c *gin.Context) {
	// get ID of URL
	id := c.Param("id")

	// Get single user
	var user []models.User
	result := initializers.DB.First(&user, id)
	if result.Error != nil{
		log.Fatal("Unable to get user:"+ result.Error.Error())
	}
	c.JSON(200, gin.H{
		"user": user,
	})
}
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	// struct to store user data
	var body struct {
		ID uint 
		Name string 
		Email *string 
		Age uint8
	}
	c.Bind(&body)

	// Get single data to update
	var user []models.User
	initializers.DB.First(&user, id)

	result := initializers.DB.Model(&user).Updates(models.User{
		ID: body.ID, 
		Name: body.Name, 
		Email: body.Email, 
		Age: body.Age})

	if result.Error != nil{
		log.Fatal("Unable to update user:"+ result.Error.Error())
	}

	c.JSON(200, gin.H{
		"user": user,
	})
}
func DeleteUser(c *gin.Context) {
	// Delete user by ID
	id := c.Param("id")

	// Delete single user
	var user []models.User
	result := initializers.DB.Delete(&user, id)
	if result.Error != nil{
		log.Fatal("Unable to delete user:"+ result.Error.Error())
	}
	c.JSON(200, gin.H{
		"deleted": id,
		"user": user,
	})
}