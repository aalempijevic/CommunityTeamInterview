package model

type Tag struct {
	Id    int    `json:"id"`
	Key   string `json:"key"`
	Label string `json:"label"`
}

type Tags []Tag
