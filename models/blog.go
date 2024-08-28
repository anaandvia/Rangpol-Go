package models

type Blog struct {
	Id    uint   `json:"id" gorm:"primary"`
	Title string `json:"title" gorm:"not null;column:title;size:254"`
	Post  string `json:"post" gorm:"not null;column:post;size:254"`
}
