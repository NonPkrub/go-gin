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
	Name string `json:"name"`
	Desc string `json:"desc"`
}

func (ca *Categories) FindAll(c *gin.Context) {
	var categories []models.Category
	ca.DB.Order("id desc").Find(&categories)

	var res []allCategoriesResponse
	copier.Copy(&res, &categories)
	c.JSON(http.StatusOK, gin.H{"categories": res})
}
func (ca *Categories) FindOne(c *gin.Context) {
	category, err := ca.findCategoryByID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err ": err.Error()})
		return
	}
	var res categoryResponse
	copier.Copy(&res, &category)
	c.JSON(http.StatusOK, gin.H{"category": res})
}

func (ca *Categories) Create(c *gin.Context) {
	var form categoryRequest
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}
	var category models.Category
	err := copier.Copy(&category, &form)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}
	if err := ca.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"category": category})
}

func (ca *Categories) Update(c *gin.Context) {
	var form updateCategory
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}
	category, err := ca.findCategoryByID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err ": err.Error()})
		return
	}

	if err := ca.DB.Model(&category).Updates(&form).Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"category": category})
}

func (ca *Categories) Delete(c *gin.Context) {
	category, err := ca.findCategoryByID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err ": err.Error()})
		return
	}
	if err := ca.DB.Unscoped().Delete(&category).Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (ca *Categories) findCategoryByID(c *gin.Context) (*models.Category, error) {
	var category models.Category
	id := c.Param("id")

	if err := ca.DB.Preload("Article").First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err ": err.Error()})
		return nil, err
	}
	return &category, nil
}
