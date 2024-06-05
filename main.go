package main

import (
	"gin/config"
	"gin/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	config.InitDB()
	defer config.CloseDB()
	//migrations.Migrate()
	//seed.Load()
	r := gin.Default()
	r.Static("/uploads", "./uploads")
	uploadDir := [...]string{"articles", "users"}
	for _, dir := range uploadDir {
		os.MkdirAll("uploads/"+dir, 0755)
	}
	routes.Serve(r)
	r.Run(":" + os.Getenv("PORT"))
}
