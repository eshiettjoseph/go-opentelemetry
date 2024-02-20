package main

import (
	"go-rest-api/controllers"
	"go-rest-api/initializers"
	"github.com/gin-gonic/gin"
)

// represents data about User data

// responds with list of all users

// Run function before main
func init(){
	// import initializers
	// initializers.LoadEnvVars()
	initializers.ConnectToDB()
}


func main(){

	router := gin.Default()


	router.POST("/api/user/", controllers.AddUser)
	router.GET("/api/users/", controllers.GetUsers)
	router.GET("/api/user/:id", controllers.GetUser)
	router.DELETE("/api/user/:id", controllers.DeleteUser)
	router.PUT("/api/user/:id", controllers.UpdateUser)
	router.PUT("/status", controllers.Status)

	router.Run()
}