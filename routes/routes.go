package routes

import (
	"gin/config"
	"gin/controllers"
	"gin/middleware"

	"github.com/gin-gonic/gin"
)

// GET /api/v1/articles
// GET /api/v1/articles/:id
// POST /api/v1/articles
// PATCH /api/v1/articles/:id
// DELETE /api/v1/articles/:id

func Serve(r *gin.Engine) {
	db := config.GetDB()
	v1 := r.Group("/api/v1")

	authGroup := v1.Group("auth")
	authController := controllers.Auth{DB: db}
	{
		authGroup.POST("register", authController.Signup)
		authGroup.POST("login", middleware.Authentication().LoginHandler)
	}

	articlesGroup := v1.Group("articles")
	articleController := controllers.Article{DB: db}
	{
		articlesGroup.GET("", articleController.FindAll)
		articlesGroup.GET("/:id", articleController.FindOne)
		articlesGroup.POST("", articleController.Create)
		articlesGroup.PATCH("/:id", articleController.Update)
		articlesGroup.DELETE("/:id ", articleController.Delete)

	}

	categoriesGroup := v1.Group("categories")
	categoriesGroupController := controllers.Categories{DB: db}
	{
		categoriesGroup.GET("", categoriesGroupController.FindAll)
		categoriesGroup.GET("/:id", categoriesGroupController.FindOne)
		categoriesGroup.POST("", categoriesGroupController.Create)
		categoriesGroup.PATCH("/:id", categoriesGroupController.Update)
		categoriesGroup.DELETE("/:id ", categoriesGroupController.Delete)
	}

}
