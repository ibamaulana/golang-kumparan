package model

import "time"

type News struct {
	ID      int64     `json:"id"`
	Author  string    `json:"author"`
	Body    string    `json:"body"`
	Created time.Time `json:"created"`
}

func (u *News) TableName() string {
	return "news"
}
