package models

import "gorm.io/gorm"
//notes model
//an entity
type Note struct { //declares a new type named Note, which is a struct
	//the json tag specifies the key name when the struct is serialized to JSON
	gorm.Model //This is an embedded struct,GORM's built-in model, includes fields ID, CreatedAt, UpdatedAt, DeletedAt
	Title   string `json:"title"`
	Content string `json:"content"`
}
