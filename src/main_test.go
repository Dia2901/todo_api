package main

import (
	"bytes"
	"encoding/json"
	"gojek/web-server-gin/pkg/cache"
	"gojek/web-server-gin/pkg/config"
	"gojek/web-server-gin/pkg/ct"
	"gojek/web-server-gin/pkg/db"
	"gojek/web-server-gin/pkg/handlers"
	"gojek/web-server-gin/pkg/models"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var (
	router = gin.Default()
)

const (
	DeleteQuery = "DELETE FROM tasks;"
	AlterQuery  = "ALTER SEQUENCE tasks_id_seq RESTART WITH 1"
)

func clearTable() {
	db.DB.Exec(DeleteQuery)
	db.DB.Exec(AlterQuery)
}

func clearCache() {
	cache.RedisClient.FlushAll()
}

func TestMain(m *testing.M) {

	config.TestInit()
	db.SetupDB()
	cache.SetupRedis()

	code := m.Run()

	clearTable()
	clearCache()
	os.Exit(code)
}

type Response struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Desc     string `json:"desc"`
	Status   bool   `json:"status"`
}

func TestAddTask(t *testing.T) {
	clearTable()
	clearCache()

	var jsonStr = []byte(`{"desc":"test product", "status": false}`)

	router.POST("/tasks/user", handlers.AddTask)
	req, _ := http.NewRequest("POST", "/tasks/user", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User", "diya")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, http.StatusCreated, w.Code, "Asserting Status Code")
	assert.Equal(t, string(responseData), "\"Task successfully Created\"")
}

func createTask(count int) {
	db.DB.Query("Insert into tasks values (1,'Testing',false,'diya')")
	cache.RedisClient.SAdd(ct.Users, "diya")
	cache.RedisClient.SAdd("diya", 1)

	if count == 1 {
		return
	}
	db.DB.Query("Insert into tasks values (2,'Testing Todo',false,'diya')")
	cache.RedisClient.SAdd(ct.Users, "diya")
	cache.RedisClient.SAdd("diya", 2)
}

func TestGetTask(t *testing.T) {
	clearTable()
	clearCache()
	createTask(1)

	router.GET("/tasks/user/:id", handlers.GetTask)
	req, _ := http.NewRequest("GET", "/tasks/user/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User", "diya")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, http.StatusOK, w.Code)

	var task models.Task
	json.Unmarshal(responseData, &task)

	assert.Equal(t, "1", task.ID)
	assert.Equal(t, "Testing", task.Desc)
	assert.Equal(t, false, task.Status)

}

func TestUpdateTask(t *testing.T) {
	clearTable()
	clearCache()
	createTask(1)

	var jsonStr = []byte(`{"status": true}`)

	router.PUT("/tasks/user/:id", handlers.UpdateTask)
	req, _ := http.NewRequest("PUT", "/tasks/user/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User", "diya")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var status bool
	assert.Equal(t, http.StatusOK, w.Code)
	rows, _ := db.DB.Query("Select status from tasks where id = 1;")
	rows.Next()
	rows.Scan(&status)

	assert.Equal(t, true, status)
}

func TestDeleteTask(t *testing.T) {
	clearTable()
	clearCache()
	createTask(1)

	router.DELETE("/tasks/user/:id", handlers.DeleteTask)
	req, _ := http.NewRequest("DELETE", "/tasks/user/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User", "diya")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var count int = 0

	db.DB.QueryRow("Select count(*) from tasks where id = 1;").Scan(&count)
	assert.Equal(t, 0, count)
}

func TestGetAllTasks(t *testing.T) {
	clearTable()
	clearCache()
	createTask(2)

	router.GET("/tasks/user", handlers.GetAllUserTasks)
	req, _ := http.NewRequest("GET", "/tasks/user", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User", "diya")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	responseData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, http.StatusOK, w.Code)
	//body := string(responseData)
	//fmt.Println(body)
	var tasks []models.Task
	json.Unmarshal(responseData, &tasks)

	assert.Equal(t, "1", tasks[0].ID)
	assert.Equal(t, "2", tasks[1].ID)
}
