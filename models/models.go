package models

type Book struct {
	ID     uint   `json:"id" gorm:"primarykey;not_null"` // auto_increment
	Title  string `json:"title"`
	Author string `json:"author"`
}
