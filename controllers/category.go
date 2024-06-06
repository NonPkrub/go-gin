package controllers

import (
	"gin/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

type Categories struct {
	DB *gorm.DB
}

type categoryResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Desc    string `json:"desc"`
	Article []struct {
		ID    uint   `json:"id"`
		Title string `json:"title"`
	} `json:"article"`
}

type categoryRequest struct {
	Name string `json:"name" binding:"required"`
	Desc string `json:"desc" binding:"required"`
}

type updateCategory struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type allCategoriesResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

func (ca *Categories) FindAll(ctx *gin.Context) {
	var categories []models.Category
	ca.DB.Order("id desc").Find(&categories)

	var res []allCategoriesResponse
	if err := copier.Copy(&res, &categories); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"categories": res})
}
func (ca *Categories) FindOne(ctx *gin.Context) {
	category, err := ca.findCategoryByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"err ": err.Error()})
		return
	}
	var res categoryResponse
	if err := copier.Copy(&res, &category); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"category": res})
}

func (ca *Categories) Create(ctx *gin.Context) {
	var form categoryRequest
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}
	var category models.Category
	err := copier.Copy(&category, &form)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}
	if err := ca.DB.Create(&category).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"category": category})
}

func (ca *Categories) Update(ctx *gin.Context) {
	var form updateCategory
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}
	category, err := ca.findCategoryByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"err ": err.Error()})
		return
	}

	if err := ca.DB.Model(&category).Updates(&form).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"category": category})
}

func (ca *Categories) Delete(ctx *gin.Context) {
	category, err := ca.findCategoryByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"err ": err.Error()})
		return
	}
	if err := ca.DB.Unscoped().Delete(&category).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (ca *Categories) findCategoryByID(ctx *gin.Context) (*models.Category, error) {
	var category models.Category
	id := ctx.Param("id")

	if err := ca.DB.Preload("Article").First(&category, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"err ": err.Error()})
		return nil, err
	}
	return &category, nil
}
