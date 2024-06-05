package migrations

import (
	"gin/models"

	"github.com/jinzhu/gorm"
	"gopkg.in/gormigrate.v1"
)

func m1717562370CreateArticleMigrationTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "1717562370",
		Migrate: func(tx *gorm.DB) error {
			err := tx.AutoMigrate(&models.Article{}).Error

			// var articles []models.Article
			// tx.Unscoped().Find(&articles)
			// for _, article := range articles {
			// 	article.CategoryID = 1
			// 	tx.Save(&article)
			// }
			return err
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.DropTable("articles").Error

			//return  tx.Model(&models.Article{}).DropColumn("category_id").Error
		},
	}
}
