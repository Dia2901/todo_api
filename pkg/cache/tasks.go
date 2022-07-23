package cache

import (
	"errors"
	"gojek/web-server-gin/pkg/ct"
	"gojek/web-server-gin/pkg/handleError"
	"gojek/web-server-gin/pkg/models"
)

func FetchTask(id string) (newTask models.Task) {
	task, _ := RedisClient.LRange(id, 0, 2).Result()
	newTask.Desc = task[0]
	newTask.ID = id
	if task[1] == ct.Done {
		newTask.Status = true
	} else {
		newTask.Status = false
	}
	return

}

func GetTask(id string, username string) (Task *models.Task, err error) {
	if RedisClient.Exists(username).Val() == 0 {
		err = errors.New("user does not exist")
		return
	}
	if !RedisClient.SIsMember(username, id).Val() {
		err = errors.New("task does not exist")
		return
	}
	if RedisClient.Exists(id).Val() == 0 {
		return
	}
	task := FetchTask(id)
	Task = &task
	return
}

func GetAllTasks(username string, tasks *[]models.Task) (notInCache []string, err error) {
	if RedisClient.Exists(username).Val() == 0 {
		err = errors.New("no such user exists")
		return
	}
	IDs, err := RedisClient.SMembers(username).Result()
	handleError.Check(err)

	for _, id := range IDs {
		if RedisClient.Exists(id).Val() == 0 {
			notInCache = append(notInCache, id)
			continue
		}
		*tasks = append(*tasks, FetchTask(id))
	}
	return

}

func AddTask(id string, username string, desc string, status bool) {
	if RedisClient.Exists(id).Val() > 0 {
		return
	}

	err := RedisClient.LPush(id, status, desc).Err()
	handleError.Check(err)

	err = RedisClient.SAdd(username, id).Err()
	handleError.Check(err)

	err = RedisClient.SAdd(ct.Users, username).Err()
	handleError.Check(err)
}

func DeleteTask(id string, username string) {
	if RedisClient.Exists(id).Val() == 0 {
		return
	}
	err := RedisClient.Del(id).Err()
	handleError.Check(err)

	err = RedisClient.SRem(username, id).Err()
	handleError.Check(err)

	if RedisClient.Exists(username).Val() == 0 {
		err = RedisClient.SRem(ct.Users, username).Err()
		handleError.Check(err)
	}
}

func UpdateTask(id string, username string, desc string, status bool) {

	if RedisClient.Exists(id).Val() != 0 {
		DeleteTask(id, username)
	}
	AddTask(id, username, desc, status)
}
