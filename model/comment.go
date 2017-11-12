package model

type Comment struct {
	Id int `json:"id"`
	//Tag Tag `json: "tag"`
	Text string	`json:"text"`
}

type Comments []Comment