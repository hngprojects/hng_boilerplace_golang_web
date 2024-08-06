package postgresql

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	defaultPage  = 1
	defaultLimit = 20
)

type Pagination struct {
	Page  int
	Limit int
}
type PaginationResponse struct {
	CurrentPage     int `json:"current_page"`
	PageCount       int `json:"page_count"`
	TotalPagesCount int `json:"total_pages_count"`
}

func GetPagination(c *gin.Context) Pagination {
	var (
		page  *int
		limit *int
	)
	if c.Query("page") != "" {
		pageInt, err := strconv.Atoi(c.Query("page"))
		if err == nil {
			page = &pageInt
		}
	}
	if c.Query("limit") != "" {
		limitInt, err := strconv.Atoi(c.Query("limit"))
		if err == nil {
			limit = &limitInt
		}
	}

	if page != nil && limit != nil {
		return Pagination{Page: *page, Limit: *limit}
	} else if page == nil && limit != nil {
		return Pagination{Page: defaultPage, Limit: *limit}
	} else if page != nil && limit == nil {
		return Pagination{Page: *page, Limit: defaultLimit}
	} else {
		return Pagination{Page: defaultPage, Limit: defaultLimit}
	}
}

func SelectAllFromDb(db *gorm.DB, order string, receiver interface{}, query interface{}, args ...interface{}) error {
	if order == "" {
		order = "desc"
	}
	tx := db.Order("id "+order).Where(query, args...).Find(receiver)
	return tx.Error
}

func SelectAllFromDbWithLimit(db *gorm.DB, order string, limit int, receiver interface{}, query interface{}, args ...interface{}) error {
	if order == "" {
		order = "desc"
	}
	tx := db.Order("id "+order).Where(query, args...).Limit(limit).Find(receiver)
	return tx.Error
}

func SelectAllFromDbOrderBy(db *gorm.DB, orderBy, order string, receiver interface{}, query interface{}, args ...interface{}) error {
	if order == "" {
		order = "desc"
	}
	if orderBy == "" {
		orderBy = "id"
	}
	tx := db.Order(orderBy+" "+order).Where(query, args...).Find(receiver)
	return tx.Error
}

func SelectAllFromByGroup(db *gorm.DB, orderBy, order string, pagination *Pagination, receiver interface{}, query interface{}, groupColumn string, args ...interface{}) (PaginationResponse, error) {

	if order == "" {
		order = "desc"
	}
	if orderBy == "" {
		orderBy = "id"
	}

	if pagination == nil {
		tx := db.Order(orderBy+" "+order).Where(query, args...).Group(groupColumn + ", id").Find(receiver)
		return PaginationResponse{}, tx.Error
	}

	var count int64
	err := db.Model(receiver).Where(query, args...).Group(groupColumn + ", id").Count(&count).Error
	if err != nil {
		return PaginationResponse{
			CurrentPage:     pagination.Page,
			PageCount:       pagination.Limit,
			TotalPagesCount: 0,
		}, err
	}

	totalPages := int(math.Ceil(float64(count) / float64(pagination.Limit)))

	tx := db.Limit(pagination.Limit).Offset((pagination.Page-1)*pagination.Limit).Order(orderBy+" "+order).Where(query, args...).Group(groupColumn + ", id").Find(receiver)
	return PaginationResponse{
		CurrentPage:     pagination.Page,
		PageCount:       int(tx.RowsAffected),
		TotalPagesCount: totalPages,
	}, tx.Error
}

func RawSelectAllFromByGroup(db *gorm.DB, orderBy, order string, pagination *Pagination, model interface{}, receiver interface{}, groupColumn string, selectQuery string, query string, args ...interface{}) (PaginationResponse, error) {

	if order == "" {
		order = "desc"
	}
	if orderBy == "" {
		orderBy = "id"
	}

	if pagination == nil {
		tx := db.Model(model).Order(orderBy+" "+order).Select(selectQuery).Where(query, args...).Group(groupColumn + ", id").Find(receiver)
		return PaginationResponse{}, tx.Error
	}

	var count int64
	err := db.Model(model).Where(query, args...).Group(groupColumn + ", id").Count(&count).Error
	if err != nil {
		return PaginationResponse{
			CurrentPage:     pagination.Page,
			PageCount:       pagination.Limit,
			TotalPagesCount: 0,
		}, err
	}

	totalPages := int(math.Ceil(float64(count) / float64(pagination.Limit)))

	tx := db.Model(model).Limit(pagination.Limit).Offset((pagination.Page-1)*pagination.Limit).Order(orderBy+" "+order).Select(selectQuery).Where(query, args...).Group(groupColumn + ", id").Find(receiver)
	return PaginationResponse{
		CurrentPage:     pagination.Page,
		PageCount:       int(tx.RowsAffected),
		TotalPagesCount: totalPages,
	}, tx.Error
}

func SelectAllFromDbOrderByPaginated(db *gorm.DB, orderBy, order string, pagination Pagination, receiver interface{}, query interface{}, args ...interface{}) (PaginationResponse, error) {

	if order == "" {
		order = "desc"
	}
	if orderBy == "" {
		orderBy = "id"
	}
	if pagination.Page <= 0 {
		pagination.Page = defaultPage
	}
	if pagination.Limit < 0 {
		pagination.Limit = defaultLimit
	}

	var count int64
	err := db.Model(receiver).Where(query, args...).Count(&count).Error
	if err != nil {
		return PaginationResponse{
			CurrentPage:     pagination.Page,
			PageCount:       pagination.Limit,
			TotalPagesCount: 0,
		}, err
	}

	totalPages := int(math.Ceil(float64(count) / float64(pagination.Limit)))

	tx := db.Limit(pagination.Limit).Offset((pagination.Page-1)*pagination.Limit).Order(orderBy+" "+order).Where(query, args...).Find(receiver)
	return PaginationResponse{
		CurrentPage:     pagination.Page,
		PageCount:       int(tx.RowsAffected),
		TotalPagesCount: totalPages,
	}, tx.Error
}

func SelectOneFromDb(db *gorm.DB, receiver interface{}, query interface{}, args ...interface{}) (error, error) {

	tx := db.Where(query, args...).First(receiver)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return tx.Error, tx.Error
	}
	return tx.Error, nil
}
func SelectLatestFromDb(db *gorm.DB, receiver interface{}, query interface{}, args ...interface{}) (error, error) {

	tx := db.Order("id desc").Where(query, args...).First(receiver)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return tx.Error, tx.Error
	}
	return tx.Error, nil
}

func SelectRandomFromDb(db *gorm.DB, receiver interface{}, query interface{}, args ...interface{}) (error, error) {

	tx := db.Order("rand()").Where(query, args...).First(receiver)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return tx.Error, tx.Error
	}
	return tx.Error, nil
}

func SelectFirstFromDb(db *gorm.DB, receiver interface{}) error {
	tx := db.First(receiver)
	return tx.Error
}

func CheckExists(db *gorm.DB, receiver interface{}, query interface{}, args ...interface{}) bool {

	tx := db.Where(query, args...).First(receiver)
	return !errors.Is(tx.Error, gorm.ErrRecordNotFound)
}

func CheckExistsInTable1(db *gorm.DB, table string, query interface{}, args ...interface{}) bool {
	var result interface{}
	tx := db.Table(table).Where(query, args...).Take(&result)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return false
		} else {
			fmt.Println("tx error", tx.Error.Error())
		}
	}

	return true
}

func CheckExistsInTable(db *gorm.DB, table string, query interface{}, args ...interface{}) bool {
	var result map[string]interface{}
	tx := db.Table(table).Where(query, args...).Take(&result)
	return tx.RowsAffected != 0
}

func PreloadEntities(db *gorm.DB, model interface{}, preloads ...string) *gorm.DB {
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	return db
}

func SelectUsersFromDb(db *gorm.DB, order string, receiver interface{}, query interface{}, args ...interface{}) error {
	if order == "" {
		order = "desc"
	}
	tx := db.Order("room_id "+order).Where(query, args...).Find(receiver)
	return tx.Error
}