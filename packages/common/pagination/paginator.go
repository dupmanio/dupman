package pagination

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Paginate(ctx *gin.Context) *Pagination {
	var err error

	pagination := &Pagination{}

	if pagination.Limit, err = strconv.Atoi(ctx.Query("limit")); err != nil {
		pagination.Limit = DefaultLimit
	}

	if pagination.Page, err = strconv.Atoi(ctx.Query("page")); err != nil {
		pagination.Page = 1
	}

	return pagination
}

func WithPagination(
	db *gorm.DB,
	value any,
	pagination *Pagination,
) func(db *gorm.DB) *gorm.DB {
	var totalItems int64

	db.Model(value).Count(&totalItems)

	pagination.TotalItems = totalItems
	pagination.TotalPages = int(math.Ceil(float64(totalItems) / float64(pagination.GetLimit())))

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}
