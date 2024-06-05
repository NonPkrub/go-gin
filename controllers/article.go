package controllers

import (
	"gin/models"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

type Article struct {
	DB *gorm.DB
}

type createdArticleResponse struct {
	ID         uint   `json:"id"`
	Title      string `json:"title"`
	Excerpt    string `json:"excerpt"`
	Body       string `json:"body"`
	Image      string `json:"image"`
	CategoryID uint   `json:"category_id"`
	category   struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
}

type createArticleRequest struct {
	Title   string                `form:"title" binding:"required"`
	Body    string                `form:"body" binding:"required"`
	Excerpt string                `form:"excerpt" binding:"required"`
	Image   *multipart.FileHeader `form:"image" binding:"required"`
}

type updateArticleRequest struct {
	Title   string                `form:"title"`
	Body    string                `form:"body"`
	Excerpt string                `form:"excerpt"`
	Image   *multipart.FileHeader `form:"image"`
}

type aresArticleResponsePaging struct {
	Item   []createdArticleResponse `json:"item"`
	Paging *pagingResult            `json:"paging"`
}

//var articles []models.Article = []models.Article{}

func (a *Article) FindAll(c *gin.Context) {
	var articles []models.Article
	//a.DB.Find(&articles)
	pagination := pagination{c: c, db: a.DB.Preload("Category").Order("id desc"), model: &articles}
	paging := pagination.pageResource()
	res := []createdArticleResponse{}
	copier.Copy(&res, &articles)

	c.JSON(http.StatusOK, gin.H{"articles": aresArticleResponsePaging{Item: res, Paging: paging}})
}

func (a *Article) FindOne(c *gin.Context) {
	article, err := a.FindArticleByID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err ": err.Error()})
		return
	}
	res := createdArticleResponse{}
	copier.Copy(&res, &article)
	c.JSON(http.StatusOK, gin.H{"article": res})
}

func (a *Article) Create(c *gin.Context) {
	var form createArticleRequest
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}

	var articles models.Article
	err := copier.Copy(&articles, &form)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}

	if err := a.DB.Create(&articles).Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}
	a.setArticle(c, &articles)
	res := createdArticleResponse{}
	copier.Copy(&res, &articles)

	c.JSON(http.StatusCreated, gin.H{"article": res})
}

func (a *Article) setArticle(c *gin.Context, article *models.Article) error {
	file, err := c.FormFile("image")
	if err != nil || file == nil {
		return err
	}

	if article.Image != "" {
		article.Image = strings.Replace(article.Image, os.Getenv("HOST"), "", 1)
	}

	pwd, _ := os.Getwd()

	os.Remove(pwd + article.Image)
	path := "uploads/articles/" + strconv.Itoa(int(article.ID))
	os.MkdirAll(path, 0755)
	fileName := path + "/" + file.Filename
	if err := c.SaveUploadedFile(file, fileName); err != nil {
		return err
	}

	article.Image = os.Getenv("HOST") + "/" + fileName
	a.DB.Save(article)

	return nil
}

func (a *Article) FindArticleByID(c *gin.Context) (*models.Article, error) {
	var article *models.Article
	id := c.Param("id")

	if err := a.DB.Preload("Category").First(&article, id).Error; err != nil {
		return nil, err
	}

	return article, nil
}

func (a *Article) Update(c *gin.Context) {
	var form updateArticleRequest
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}

	article, err := a.FindArticleByID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err ": err.Error()})
		return
	}

	if err := a.DB.Model(&article).Updates(&form).Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}

	a.setArticle(c, article)

	res := createdArticleResponse{}
	copier.Copy(&res, &article)
	c.JSON(http.StatusOK, gin.H{"article": res})

}

func (a *Article) Delete(c *gin.Context) {
	article, err := a.FindArticleByID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"err ": err.Error})
	}

	a.DB.Unscoped().Delete(&article)
	c.Status(http.StatusNoContent)
}
