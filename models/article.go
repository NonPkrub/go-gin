package models

import "github.com/jinzhu/gorm"

type Article struct {
	gorm.Model
	Title      string `gorm:"unique;not null" json:"title"`
	Excerpt    string `gorm:"unique;not null" json:"excerpt"`
	Body       string `gorm:"not null" json:"body"`
	Image      string `gorm:"not null" json:"image"`
	CategoryID uint
	Category   Category
}
