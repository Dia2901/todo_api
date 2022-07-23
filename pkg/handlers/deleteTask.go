package handlers

import (
	"fmt"
	"gojek/web-server-gin/pkg/cache"
	"gojek/web-server-gin/pkg/ct"
	"gojek/web-server-gin/pkg/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteTask(c *gin.Context) {
	id := c.Param(ct.ID)
	user := c.Request.Header[ct.User]
	username := user[0]

	err := db.DeleteTask(id, username)
	if err != nil {
		fmt.Println(err)
		c.IndentedJSON(http.StatusNotFound, "Task not found")
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Task successfully  deleted"})
	cache.DeleteTask(id, username)

}
