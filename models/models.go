package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Post struct {
	gorm.Model
	Title string `json:"title"`
	Body string `json:"body"`
	UserID uint `json:"user_id"`
}

