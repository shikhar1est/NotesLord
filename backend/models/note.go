package models

import "gorm.io/gorm"
//notes model
//an entity
type Note struct {
	gorm.Model
	Title   string `json:"title"`
	Content string `json:"content"`
}
