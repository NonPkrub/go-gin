package controllers

import (
	"gin/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

type Auth struct {
	DB *gorm.DB
}

type authForm struct {
	Email    string `form:"email" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type authResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
}

func (a *Auth) Signup(c *gin.Context) {
	var form authForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}

	var user models.User
	copier.Copy(&user, &form)
	user.Password = user.GeneratePassword()
	if err := a.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}
	res := authResponse{}
	copier.Copy(&res, &user)
	c.JSON(http.StatusCreated, gin.H{"user": res})
}

// func (a *Auth) Login(c *gin.Context) {}
