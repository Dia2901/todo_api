package handlers

import (
	"fmt"
	"gojek/web-server-gin/pkg/cache"
	"gojek/web-server-gin/pkg/ct"
	"gojek/web-server-gin/pkg/db"
	"gojek/web-server-gin/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllUserTasks(c *gin.Context) {

	user := c.Request.Header[ct.User]
	username := user[0]
	var notInCache []string
	var tasks []models.Task

	notInCache, err := cache.GetAllTasks(username, &tasks)

	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusNotFound, "user not found")
		return
	}
	tasks = append(tasks, db.GetAllTasks(notInCache, username)...)
	c.IndentedJSON(http.StatusOK, tasks)

}
