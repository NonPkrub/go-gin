package controllers

import (
	"gin/config"
	"gin/models"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

type User struct {
	DB *gorm.DB
}

type userResponse struct {
	ID     uint   `json:"id"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
	Role   string `json:"role"`
}

type userCreateRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

type userUpdateRequest struct {
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty, min=8"`
	Name     string `json:"name"`
}

type userPaging struct {
	Items  []userResponse `json:"items"`
	Paging *pagingResult  `json:"paging"`
}

func setUserImage(ctx *gin.Context, user *models.User) error {
	file, _ := ctx.FormFile("avatar")
	if file == nil {
		return nil
	}

	if user.Avatar != "" {
		user.Avatar = strings.Replace(user.Avatar, os.Getenv("HOST"), "", 1)
		pwd, _ := os.Getwd()
		os.Remove(pwd + user.Avatar)
	}

	path := "uploads/users/" + strconv.Itoa(int(user.ID))
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}
	fileName := path + "/" + file.Filename
	if err := ctx.SaveUploadedFile(file, fileName); err != nil {
		return err
	}

	db := config.GetDB()
	user.Avatar = os.Getenv("HOST") + "/" + fileName
	db.Save(user)
	return nil
}

func (u *User) Create(ctx *gin.Context) {

	var form userCreateRequest
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}

	var user models.User
	_ = copier.Copy(&user, &form)
	user.Password = user.GeneratePassword()
	if err := u.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}

	res := userResponse{}
	_ = copier.Copy(&res, &user)
	ctx.JSON(http.StatusCreated, gin.H{"user": res})
}

// api/v1/users?term=as
func (u *User) FindAll(ctx *gin.Context) {
	var users []models.User

	query := u.DB.Order("id desc").Find(&users)
	term := ctx.Query("term")
	if term != "" {
		query = query.Where("name ILIKE ?", "%"+term+"%")
	}

	pagination := pagination{ctx: ctx, db: query, model: &users}
	paging := pagination.pageResource()

	//var res []userResponse //nil slice
	res := []userResponse{} //empty slice
	_ = copier.Copy(&res, &users)
	ctx.JSON(http.StatusOK, gin.H{"users": userPaging{Items: res, Paging: paging}})
}

func (u *User) FindOne(ctx *gin.Context) {

	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"err ": err.Error()})
		return
	}

	var res userResponse
	_ = copier.Copy(&res, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": res})
}

func (u *User) findUserByID(ctx *gin.Context) (*models.User, error) {

	id := ctx.Param("id")
	var user models.User

	if err := u.DB.First(&user, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"err ": err.Error()})
		return nil, err
	}

	return &user, nil
}

func (u *User) Update(ctx *gin.Context) {

	var form userUpdateRequest
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}

	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"err ": err.Error()})
		return
	}

	if form.Password != "" {
		user.Password = user.GeneratePassword()
	}

	if err := u.DB.Model(&user).Updates(&form).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}

	var res userResponse
	_ = copier.Copy(&res, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": res})
}

func (u *User) Delete(ctx *gin.Context) {

	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"err ": err.Error()})
		return
	}

	if err := u.DB.Unscoped().Delete(&user).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (u *User) Promote(ctx *gin.Context) {

	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"err ": err.Error()})
		return
	}

	user.Promote()
	if err := u.DB.Save(&user).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}

	var res userResponse
	_ = copier.Copy(&res, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": res})
}

func (u *User) Demote(ctx *gin.Context) {
	//authorization
	// sub, _ := ctx.Get("sub")
	// if sub.(*models.User).Role != "admin" {
	// 	ctx.JSON(http.StatusForbidden, gin.H{"err ": "Forbidden"})
	// 	return
	// }

	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"err ": err.Error()})
		return
	}

	user.Demote()
	if err := u.DB.Save(&user).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}

	var res userResponse
	_ = copier.Copy(&res, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": res})
}
