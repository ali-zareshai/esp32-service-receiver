package utils

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

func Paginate(req *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		pageSize, err := strconv.Atoi(req.DefaultQuery("size", "100"))
		if err != nil {
			return db
		}
		pageNumber, err := strconv.Atoi(req.DefaultQuery("page", "1"))
		if err != nil {
			return db
		}
		if pageNumber < 1 {
			pageNumber = 1
		}
		offset := (pageNumber - 1) * pageSize

		return db.Order("created_at").Limit(pageSize).Offset(offset)
	}
}
