package model

type Comment struct {
	Id int `json:"id"`

	Text string `json:"text"`
}

type Comments []Comment
