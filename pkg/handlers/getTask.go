package handlers

import (
	"fmt"
	"gojek/web-server-gin/pkg/cache"
	"gojek/web-server-gin/pkg/ct"
	"gojek/web-server-gin/pkg/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTask(c *gin.Context) {
	id := c.Param(ct.ID)
	User := c.Request.Header[ct.User]
	username := User[0]

	task, err := cache.GetTask(id, username)

	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusForbidden, "Access Denied")
		return
	}
	if task != nil {
		c.IndentedJSON(http.StatusOK, *task)
		return
	}
	tasks, err := db.GetTask(id, username)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusNotFound, "task not found")
		return
	}
	c.IndentedJSON(http.StatusOK, *tasks)
	Task := *tasks
	cache.AddTask(Task.ID, username, Task.Desc, Task.Status)
}
