package gglmm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// GormOpenConfig --
func GormOpenConfig(config ConfigDB) *gorm.DB {
	if !config.Check() {
		log.Printf("%+v\n", config)
		log.Fatal("DBConfig invalid")
	}
	connMaxLifetime, err := time.ParseDuration(fmt.Sprintf("%ds", config.ConnMaxLifetime))
	if err != nil {
		log.Fatal(err)
	}
	url := config.User + ":" + config.Password + "@(" + config.Address + ")/" + config.Database + "?charset=utf8mb4&parseTime=true&loc=UTC"
	return GormOpen(
		config.Dialect,
		url,
		config.MaxOpen,
		config.MaxIdel,
		connMaxLifetime,
	)
}

// GormOpen --
func GormOpen(dialect string, url string, maxOpen int, maxIdle int, connMaxLifetime time.Duration) *gorm.DB {
	var err error
	db, err := gorm.Open(dialect, url)
	if err != nil {
		panic(err)
	}

	db.SingularTable(true)

	sqlDB := db.DB()
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	return db
}

func gormFilterRequest(db *gorm.DB, filterRequest FilterRequest) (*gorm.DB, error) {
	db, err := gormFilters(db, filterRequest.Filters)
	if err != nil {
		return nil, err
	}
	if filterRequest.Order != "" {
		db = db.Order(filterRequest.Order)
	}
	return db, nil
}

func gormFilters(db *gorm.DB, filters []Filter) (*gorm.DB, error) {
	var err error
	for _, filter := range filters {
		db, err = gormFilter(db, filter)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

func gormFilter(db *gorm.DB, filter Filter) (*gorm.DB, error) {
	if !filter.Check() {
		return nil, ErrFilter
	}
	if filter.Field == FilterFieldDeleted {
		if deleted, ok := filter.Value.(string); ok {
			if deleted == FilterValueAll.Value {
				return db.Unscoped(), nil
			}
			if deleted == FilterValueDeleted.Value {
				return db.Unscoped().Where("deleted_at is not null"), nil
			}
		}
		return db, nil
	}
	switch filter.Operate {
	case FilterOperateEqual:
		return db.Where(filter.Field+" = ?", filter.Value), nil
	case FilterOperateNotEqual:
		return db.Where(filter.Field+" <> ?", filter.Value), nil
	case FilterOperateGreaterThan:
		return db.Where(filter.Field+" > ?", filter.Value), nil
	case FilterOperateGreaterEqual:
		return db.Where(filter.Field+" >= ?", filter.Value), nil
	case FilterOperateLessThan:
		return db.Where(filter.Field+" < ?", filter.Value), nil
	case FilterOperateLessEqual:
		return db.Where(filter.Field+" <= ?", filter.Value), nil
	case FilterOperateIn:
		return gormFilterIn(db, filter)
	case FilterOperateBetween:
		return gormFilterBetween(db, filter)
	case FilterOperateLike:
		return gormFilterLike(db, filter)
	default:
		return nil, ErrFilterOperate
	}
}

func gormFilterIn(db *gorm.DB, filter Filter) (*gorm.DB, error) {
	values, ok := filter.Value.([]interface{})
	if !ok {
		return nil, ErrFilterValueType
	}
	return db.Where(filter.Field+" in (?)", values), nil
}

func gormFilterBetween(db *gorm.DB, filter Filter) (*gorm.DB, error) {
	values, ok := filter.Value.([]interface{})
	if !ok {
		return nil, ErrFilterValueType
	}
	if len(values) != 2 {
		return nil, ErrFilterValueSize
	}
	if values[0] != nil && values[1] != nil {
		return db.Where(filter.Field+" between ? and ?", values[0], values[1]), nil
	} else if values[0] != nil && values[1] == nil {
		return db.Where(filter.Field+" >= ?", values[0]), nil
	} else if values[0] == nil && values[1] != nil {
		return db.Where(filter.Field+" <= ?", values[1]), nil
	} else {
		return db, nil
	}
}

func gormFilterLike(db *gorm.DB, filter Filter) (*gorm.DB, error) {
	stringValue, ok := filter.Value.(string)
	if !ok {
		return nil, ErrFilterValueType
	}
	fields := strings.Split(filter.Field, FilterSeparator)
	values := strings.Split(stringValue, FilterSeparator)
	var wheres = []string{}
	var likes []interface{}
	count := 0
	wheres = append(wheres, "(")
	for _, field := range fields {
		for _, value := range values {
			if count == 0 {
				wheres = append(wheres, field, "like", "?")
			} else {
				wheres = append(wheres, "or", field, "like", "?")
			}
			likes = append(likes, "%"+value+"%")
			count++
		}
	}
	wheres = append(wheres, ")")
	return db.Where(strings.Join(wheres, " "), likes...), nil
}
