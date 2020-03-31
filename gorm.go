package gglmm

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// NewGormDBConfig --
func NewGormDBConfig(config ConfigDB) *gorm.DB {
	if !config.Check() {
		log.Printf("%+v\n", config)
		log.Fatal("DBConfig invalid")
	}
	connMaxLifetime, err := time.ParseDuration(fmt.Sprintf("%ds", config.ConnMaxLifetime))
	if err != nil {
		log.Fatal(err)
	}
	url := config.User + ":" + config.Password + "@(" + config.Address + ")/" + config.Database + "?charset=utf8mb4&parseTime=true&loc=UTC"
	return NewGormDB(
		config.Dialect,
		url,
		config.MaxOpen,
		config.MaxIdel,
		connMaxLifetime,
	)
}

// NewGormDB --
func NewGormDB(dialect string, url string, maxOpen int, maxIdle int, connMaxLifetime time.Duration) *gorm.DB {
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

var (
	// ErrGormRecordNotFound --
	ErrGormRecordNotFound = gorm.ErrRecordNotFound
	// ErrFilterValueType --
	ErrFilterValueType = errors.New("过滤值类型错误")
	// ErrFilterValueSize --
	ErrFilterValueSize = errors.New("过滤值大小错误")
	// ErrFilterOperate --
	ErrFilterOperate = errors.New("过滤操作错误")
)

func gormSetupFilterRequest(db *gorm.DB, filterRequest FilterRequest) (*gorm.DB, error) {
	db, err := gormSetupFilters(db, filterRequest.Filters)
	if err != nil {
		return nil, err
	}
	if filterRequest.Order != "" {
		db = db.Order(filterRequest.Order)
	}
	return db, nil
}

func gormSetupFilters(db *gorm.DB, filters []Filter) (*gorm.DB, error) {
	var err error
	for _, filter := range filters {
		db, err = gormSetupFilter(db, filter)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

func gormSetupFilter(db *gorm.DB, filter Filter) (*gorm.DB, error) {
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
		return gormSetupFilterIn(db, filter)
	case FilterOperateBetween:
		return gormSetupFilterBetween(db, filter)
	case FilterOperateLike:
		return gormSetupFilterLike(db, filter)
	}
	return nil, ErrFilterOperate
}

func gormSetupFilterIn(db *gorm.DB, filter Filter) (*gorm.DB, error) {
	stringValue, ok := filter.Value.(string)
	if !ok {
		return nil, ErrFilterValueType
	}
	values := strings.Split(stringValue, FilterSeparator)
	return db.Where(filter.Field+" in (?)", values), nil
}

func gormSetupFilterBetween(db *gorm.DB, filter Filter) (*gorm.DB, error) {
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

func gormSetupFilterLike(db *gorm.DB, filter Filter) (*gorm.DB, error) {
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
