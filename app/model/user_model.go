package model

type User struct {
	Id        string   `json:"id"`
	Interests []string `json:"interests"`
}
