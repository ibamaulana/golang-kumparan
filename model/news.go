package model

import "time"

type News struct {
	ID      int       `json:"id"`
	Author  string    `json:"author"`
	Body    string    `json:"body"`
	Created time.Time `json:"created"`
}

type ElasticRequest struct {
	Created time.Time `json:"created"`
}

type ElasticResponse struct {
	ID      int       `json:"id"`
	Created time.Time `json:"created"`
}

func (u *News) TableName() string {
	return "news"
}
