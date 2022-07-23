package models

type Task struct {
	ID     string `json:"id"`
	Desc   string `json:"desc"`
	Status bool   `json:"status"`
}
