package main

import (
	"gojek/web-server-gin/pkg/cache"
	"gojek/web-server-gin/pkg/config"
	"gojek/web-server-gin/pkg/db"
	"gojek/web-server-gin/pkg/handlers"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {

	config.Init()
	db.SetupDB()
	cache.SetupRedis()

	router := gin.Default()
	router.GET("/tasks/user", handlers.GetAllUserTasks)
	router.GET("/tasks/user/:id", handlers.GetTask)
	router.POST("/tasks/user", handlers.AddTask)
	router.DELETE("/tasks/user/:id", handlers.DeleteTask)
	router.PUT("/tasks/user/:id", handlers.UpdateTask)
	router.Run("localhost:8000")
}
