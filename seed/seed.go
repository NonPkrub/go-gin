package seed

import (
	"gin/config"
	"gin/migrations"
	"gin/models"
	"log"
	"math/rand"
	"strconv"

	"github.com/go-faker/faker/v4"
)

func Load() {
	db := config.GetDB()

	db.DropTableIfExists("articles", "categories", "migrations")
	migrations.Migrate()
	log.Println("Creating categories...")

	numOfCategories := 20
	categories := make([]models.Category, 0, numOfCategories)
	for i := 1; i <= numOfCategories; i++ {
		category := models.Category{
			Name: faker.Word(),
			Desc: faker.Paragraph(),
		}

		db.Create(&category)
		categories = append(categories, category)

	}

	log.Println("Creating categories...")

	numOfArticle := 50
	articles := make([]models.Article, 0, numOfArticle)
	for i := 1; i <= numOfArticle; i++ {
		article := models.Article{
			Title:      faker.Sentence(),
			Excerpt:    faker.Sentence(),
			Body:       faker.Paragraph(),
			Image:      "https://source.unsplash.com/random/300x200?sig=" + strconv.Itoa(i),
			CategoryID: uint(rand.Intn(numOfCategories) + 1),
		}

		db.Create(&article)
		articles = append(articles, article)
	}
}
