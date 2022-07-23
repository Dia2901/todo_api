package db

import (
	"bytes"
	"errors"
	"fmt"
	"gojek/web-server-gin/pkg/cache"
	"gojek/web-server-gin/pkg/handleError"
	"gojek/web-server-gin/pkg/models"
)

const (
	SelectQuery        = "SELECT id from tasks where username = $1 and des = $2"
	CountQuery         = "SELECT count(*) from tasks where id = $1 and username = $2"
	GetQuery           = "Select * from tasks where id = $1 and username = $2"
	InsertQuery        = "INSERT INTO tasks( username,des,status) values ( $1,$2,$3) Returning id"
	DeleteQuery        = "DELETE FROM tasks where id = $1 and username = $2"
	UpdateQuery        = "UPDATE tasks SET status = $1 where id = $2 and username = $3 returning des"
	tableCreationQuery = `CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		des TEXT,
		status bool,
		username TEXT
	)`
)

func createTable() {
	_, err := DB.Exec(tableCreationQuery)
	handleError.Check(err)
}

func GetTask(id string, username string) (task *models.Task, err error) {

	rows, err := DB.Query(GetQuery, id, username)
	handleError.Check(err)
	var ID string
	var desc string
	var status bool

	rows.Next()
	rows.Scan(&ID, &desc, &status, &username)
	if ID == "" {
		err = errors.New("task not found in db")
		return
	}
	task = &models.Task{ID: id, Desc: desc, Status: status}
	return
}

func CreateQuery(notInCache []string) (Query string) {
	buf := bytes.NewBufferString("SELECT * FROM tasks where id in (")
	for i, v := range notInCache {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(v)
	}
	buf.WriteString(");")

	Query = buf.String()
	return
}

func GetAllTasks(notInCache []string, username string) (dbtasks []models.Task) {
	if len(notInCache) == 0 {
		return
	}

	rows, err := DB.Query(CreateQuery(notInCache))
	handleError.Check(err)

	for rows.Next() {
		var id string
		var username string
		var desc string
		var status bool

		rows.Scan(&id, &desc, &status, &username)
		if id != "" {
			cache.AddTask(id, username, desc, status)
		}

		dbtasks = append(dbtasks, models.Task{ID: id, Desc: desc, Status: status})
	}

	return dbtasks

}

func AddTask(newTask models.Task, username string) (id string, err error) {

	rows, _ := DB.Query(SelectQuery, username, newTask.Desc)

	rows.Next()
	rows.Scan(&id)
	fmt.Println(rows)
	if id != "" {
		err = errors.New("task already exists")
		return
	}

	rows, er := DB.Query(InsertQuery, username, newTask.Desc, newTask.Status)
	handleError.Check(er)
	rows.Next()
	rows.Scan(&id)
	return
}

func DeleteTask(id string, username string) (err error) {

	var count int = 0
	DB.QueryRow(CountQuery, id, username).Scan(&count)

	if count == 0 {
		err = errors.New("task does not exist")
		return
	}
	_, err = DB.Query(DeleteQuery, id, username)
	handleError.Check(err)
	return
}

func UpdateTask(id string, username string, status bool) (des string, err error) {

	var count int = 0

	DB.QueryRow(CountQuery, id, username).Scan(&count)
	if count == 0 {
		err = errors.New("task does not exist ")
		des = ""
		return
	}

	rows, err := DB.Query(UpdateQuery, status, id, username)
	handleError.Check(err)

	rows.Next()
	rows.Scan(&des)
	return

}
