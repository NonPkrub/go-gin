package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Dashboard struct {
	DB *gorm.DB
}

type dashboardArticle struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Excerpt string `json:"excerpt"`
	Image   string `json:"image"`
}

type dashboardResponse struct {
	LatestArticles []dashboardArticle `json:"latest_articles"`
	UserCount      []struct {
		Role  string `json:"role"`
		Count string `json:"count"`
	} `json:"user_count"`
	CategoryCount uint `json:"category_count"`
	ArticleCount  uint `json:"articale_count"`
}

func (d *Dashboard) GetDashboard(ctx *gin.Context) {
	res := dashboardResponse{}
	d.DB.Table("articles").Order("id desc").Limit(5).Find(&res.LatestArticles)
	d.DB.Table("articles").Count(&res.ArticleCount)
	d.DB.Table("categories").Count(&res.CategoryCount)
	d.DB.Table("users").Select("role, count(*) as count").Group("role").Scan(&res.UserCount)

	ctx.JSON(http.StatusOK, gin.H{"dashboard": &res})
}
