package controllers

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type pagingResult struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	PrevPage   int `json:"prev_page"`
	NextPage   int `json:"next_page"`
	Count      int `json:"count"`
	TotalPages int `json:"total_pages"`
}

type pagination struct {
	c     *gin.Context
	db    *gorm.DB
	model interface{}
}

func (p *pagination) pageResource() *pagingResult {
	page, _ := strconv.Atoi(p.c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(p.c.DefaultQuery("limit", "10"))

	// var count int
	// db.Model(model).Count(&count)
	ch := make(chan int)
	go p.countRecord(ch)

	offset := (page - 1) * limit
	p.db.Limit(limit).Offset(offset).Find(p.model)

	count := <-ch
	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	var nextPage int
	if page == totalPages {
		nextPage = totalPages
	} else {
		nextPage = page - 1
	}

	return &pagingResult{
		Page:       page,
		Limit:      limit,
		Count:      count,
		PrevPage:   page - 1,
		NextPage:   nextPage,
		TotalPages: totalPages,
	}

}

func (p *pagination) countRecord(ch chan int) {
	var count int
	p.db.Model(p.model).Count(&count)

	ch <- count
}
