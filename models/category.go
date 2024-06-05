package models

import "github.com/jinzhu/gorm"

type Category struct {
	gorm.Model
	Name    string `gorm:"unique;not null" json:"name"`
	Desc    string `gorm:"not null" json:"desc"`
	Article []Article
}
