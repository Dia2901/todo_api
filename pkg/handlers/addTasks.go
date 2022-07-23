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

func AddTask(c *gin.Context) {
	var newTask models.Task

	err := c.BindJSON(&newTask)
	handleError.Check(err)

	user := c.Request.Header[ct.User]
	username := user[0]

	id, err := db.AddTask(newTask, username)

	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusConflict, "Task Already Exists")
	} else {
		c.IndentedJSON(http.StatusCreated, "Task successfully Created")
	}

	cache.AddTask(id, username, newTask.Desc, newTask.Status)

}
