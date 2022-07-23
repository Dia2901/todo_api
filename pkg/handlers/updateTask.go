package handlers

import (
	"fmt"
	"gojek/web-server-gin/pkg/cache"
	"gojek/web-server-gin/pkg/ct"
	"gojek/web-server-gin/pkg/db"
	"gojek/web-server-gin/pkg/handleError"
	"gojek/web-server-gin/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateTask(c *gin.Context) {

	id := c.Param(ct.ID)
	user := c.Request.Header[ct.User]
	username := user[0]

	var newTask models.Task

	err := c.BindJSON(&newTask)
	handleError.Check(err)

	desc, err := db.UpdateTask(id, username, newTask.Status)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusNotFound, "task not found")
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Task successfully Updated"})
	cache.UpdateTask(id, username, desc, newTask.Status)

}
