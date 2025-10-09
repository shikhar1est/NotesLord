package models

import "gorm.io/gorm"
//notes model
type Note struct {
	gorm.Model
	Title   string `json:"title"`
	Content string `json:"content"`
}
