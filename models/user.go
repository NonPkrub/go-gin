package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Name     string `gorm:"unique;not null" json:"name"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"password"`
	Avatar   string `"json:"avatar"`
	Role     string `"gorm:"defualt:'Member';not null" json:"role"`
}

func (u *User) GeneratePassword() string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	return string(hash)
}

func (u *User) Promote() {
	u.Role = "Editor"
}

func (u *User) Demote() {
	u.Role = "Member"
}
