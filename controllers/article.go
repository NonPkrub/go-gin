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
	Category   struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
	User struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	} `json:"user"`
}

type createOrUpdateArticleResponse struct {
	ID         uint   `json:"id"`
	Title      string `json:"title"`
	Excerpt    string `json:"excerpt"`
	Body       string `json:"body"`
	Image      string `json:"image"`
	CategoryID uint   `json:"category_id"`
	UserID     uint   `json:"user_id"`
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

// var articles []models.Article = []models.Article{}
// api/v1/articles?category_id=1&term=go
func (a *Article) FindAll(ctx *gin.Context) {
	articles := []models.Article{} //empty slice

	query := a.DB.Preload("User").Preload("Category").Order("id desc")

	categoriesID := ctx.Query("category_id")
	term := ctx.Query("term")
	if categoriesID != "" {
		query = query.Where("category_id = ?", categoriesID)
	}
	if term != "" {
		query = query.Where("title ILIKE?", "%"+term+"%")
	}

	//a.DB.Find(&articles)
	pagination := pagination{ctx: ctx, db: query, model: &articles}
	paging := pagination.pageResource()
	//var res []createdArticleResponse // nil slice
	res := []createdArticleResponse{} // empty slice
	if err := copier.Copy(&res, &articles); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"articles": aresArticleResponsePaging{Item: res, Paging: paging}})
}

func (a *Article) FindOne(ctx *gin.Context) {
	article, err := a.FindArticleByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"err ": err.Error()})
		return
	}
	res := createdArticleResponse{}
	if err := copier.Copy(&res, &article); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"article": res})
}

func (a *Article) Create(ctx *gin.Context) {
	var form createArticleRequest
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}

	var article models.Article
	user, _ := ctx.Get("sub")
	if err := copier.Copy(&article, user); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}
	article.User = *user.(*models.User)

	var articles models.Article
	_ = copier.Copy(&articles, &form)
	if err := a.DB.Create(&articles).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}
	if err := a.setArticle(ctx, &articles); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}
	res := createOrUpdateArticleResponse{}
	_ = copier.Copy(&res, &articles)
	ctx.JSON(http.StatusCreated, gin.H{"article": res})
}

func (a *Article) setArticle(ctx *gin.Context, article *models.Article) error {
	file, err := ctx.FormFile("image")
	if err != nil || file == nil {
		return err
	}

	if article.Image != "" {
		article.Image = strings.Replace(article.Image, os.Getenv("HOST"), "", 1)
	}

	pwd, _ := os.Getwd()

	os.Remove(pwd + article.Image)
	path := "uploads/articles/" + strconv.Itoa(int(article.ID))
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	fileName := path + "/" + file.Filename
	if err := ctx.SaveUploadedFile(file, fileName); err != nil {
		return err
	}

	article.Image = os.Getenv("HOST") + "/" + fileName
	a.DB.Save(article)

	return nil
}

func (a *Article) FindArticleByID(ctx *gin.Context) (*models.Article, error) {
	var article *models.Article
	id := ctx.Param("id")

	if err := a.DB.Preload("User").Preload("Category").First(&article, id).Error; err != nil {
		return nil, err
	}

	return article, nil
}

func (a *Article) Update(ctx *gin.Context) {
	var form updateArticleRequest
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}

	article, err := a.FindArticleByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"err ": err.Error()})
		return
	}

	if err := a.DB.Model(&article).Updates(&form).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}

	if err := a.setArticle(ctx, article); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}

	res := createOrUpdateArticleResponse{}
	if err := copier.Copy(&res, &article); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"err ": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"article": res})

}

func (a *Article) Delete(ctx *gin.Context) {
	article, err := a.FindArticleByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"err ": err.Error})
	}

	a.DB.Unscoped().Delete(&article)
	ctx.Status(http.StatusNoContent)
}
