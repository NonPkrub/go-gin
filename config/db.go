package config

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func InitDB() {
	var err error
	// db, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	db, err := gorm.Open("postgres", os.Getenv("DB_CONNECT"))
	if err != nil {
		log.Fatal(err)
	}

	db.LogMode(gin.Mode() == gin.DebugMode)
	// db.AutoMigrate(&models.Article{})
}

func GetDB() *gorm.DB {
	return db
}
func CloseDB() {
	db.Close()
}
