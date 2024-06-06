package controllers

import (
	"gin/models"
	"mime/multipart"
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

type updateProfileRequest struct {
	Email  string                `form:"email"`
	Name   string                `form:"name"`
	Avatar *multipart.FileHeader `form:"avatar"`
}

func (a *Auth) GetMe(ctx *gin.Context) {
	sub, _ := ctx.Get("sub")
	user := sub.(*models.User)

	var myProfile userResponse
	if err := copier.Copy(&myProfile, &user); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"user": myProfile})
}

func (a *Auth) Signup(ctx *gin.Context) {
	var form authForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}

	var user models.User
	if err := copier.Copy(&user, &form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}
	user.Password = user.GeneratePassword()
	if err := a.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}
	res := authResponse{}
	if err := copier.Copy(&res, &user); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"user": res})
}

func (a *Auth) UpdateProfile(ctx *gin.Context) {
	sub, _ := ctx.Get("sub")
	user := sub.(*models.User)

	var form updateProfileRequest
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}

	if err := setUserImage(ctx, user); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}
	if err := a.DB.Model(&user).Update(&form).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error})
		return
	}
	var myProfile userResponse
	if err := copier.Copy(&myProfile, &user); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"user": myProfile})
}
