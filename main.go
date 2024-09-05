package main

import (
	"fmt"
	"os"
	"verba/config"
	"verba/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()
	defer config.DB.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()

	r.POST("/tasks", handlers.CreateTask)
	r.GET("/tasks", handlers.GetAllTasks)
	r.GET("/tasks/:id", handlers.GetTask)
	r.PUT("/tasks/:id", handlers.UpdateTask)
	r.DELETE("/tasks/:id", handlers.DeleteTask)
	r.Run()
	fmt.Printf("Server is running on port %s", port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
