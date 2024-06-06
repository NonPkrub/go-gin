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

	authenticate := middleware.Authentication().MiddlewareFunc()
	authorization := middleware.Authorize()
	authGroup := v1.Group("auth")
	authController := controllers.Auth{DB: db}
	{
		authGroup.POST("/register", authController.Signup)
		authGroup.POST("/login", middleware.Authentication().LoginHandler)
		authGroup.GET("/profile", authenticate, authController.GetMe)
		authGroup.PATCH("/profile", authenticate, authController.UpdateProfile)
	}

	articlesGroup := v1.Group("articles")
	articleController := controllers.Article{DB: db}
	articlesGroup.GET("", articleController.FindAll)
	articlesGroup.GET("/:id", articleController.FindOne)
	articlesGroup.Use(authenticate, authorization)
	{
		articlesGroup.POST("", authenticate, articleController.Create)
		articlesGroup.PATCH("/:id", articleController.Update)
		articlesGroup.DELETE("/:id ", articleController.Delete)

	}

	categoriesGroup := v1.Group("categories")
	categoriesGroup.Use(authenticate, authorization)
	categoriesGroupController := controllers.Categories{DB: db}
	{
		categoriesGroup.GET("", categoriesGroupController.FindAll)
		categoriesGroup.GET("/:id", categoriesGroupController.FindOne)
		categoriesGroup.POST("", categoriesGroupController.Create)
		categoriesGroup.PATCH("/:id", categoriesGroupController.Update)
		categoriesGroup.DELETE("/:id ", categoriesGroupController.Delete)
	}

	userGroup := v1.Group("users")
	userGroup.Use(authenticate, authorization)
	userGroupController := controllers.User{DB: db}
	{
		userGroup.GET("", userGroupController.FindAll)
		userGroup.POST("", userGroupController.Create)
		userGroup.GET("/:id", userGroupController.FindOne)
		userGroup.PATCH("/:id", userGroupController.Update)
		userGroup.DELETE("/:id ", userGroupController.Delete)
		userGroup.PATCH("/:id/promote", userGroupController.Promote)
		userGroup.PATCH("/:id/demote", userGroupController.Demote)

	}

	dashGroup := v1.Group("dashboard")
	dashGroupController := controllers.Dashboard{DB: db}
	dashGroup.Use(authenticate, authorization)
	{
		dashGroup.GET("", dashGroupController.GetDashboard)
	}
}
