package gglmm

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

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

var gormDB *gorm.DB = nil

// RegisterGormDB 注册全局GormDB
func RegisterGormDB(dialect string, url string, maxOpen int, maxIdle int, connMaxLifetime time.Duration) {
	gormDB = NewGormDB(dialect, url, maxOpen, maxIdle, connMaxLifetime)
}

// RegisterGormDBConfig --
func RegisterGormDBConfig(config ConfigDB) {
	if !config.Check() {
		log.Printf("%+v\n", config)
		log.Fatal("DBConfig invalid")
	}
	connMaxLifetime, err := time.ParseDuration(fmt.Sprintf("%ds", config.ConnMaxLifetime))
	if err != nil {
		log.Fatal(err)
	}
	url := config.User + ":" + config.Password + "@(" + config.Address + ")/" + config.Database + "?charset=utf8mb4&parseTime=true&loc=UTC"
	gormDB = NewGormDB(
		config.Dialect,
		url,
		config.MaxOpen,
		config.MaxIdel,
		connMaxLifetime,
	)
}

// NewGormDB --
func NewGormDB(dialect string, url string, maxOpen int, maxIdle int, connMaxLifetime time.Duration) *gorm.DB {
	gormDB, err := gorm.Open(dialect, url)
	if err != nil {
		panic(err)
	}

	gormDB.SingularTable(true)

	sqlDB := gormDB.DB()
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	return gormDB
}

// CloseGormDB 关闭全局数据库连接
func CloseGormDB() {
	if gormDB != nil {
		gormDB.Close()
	}
}

// GormDB --
func GormDB() *gorm.DB {
	return gormDB
}

// GormBegin --
func GormBegin() *gorm.DB {
	return gormDB.Begin()
}

// GormNewRecord --
func GormNewRecord(model interface{}) bool {
	return gormDB.NewRecord(model)
}

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
	var where = "("
	var likes []interface{}
	for index, field := range fields {
		for _, value := range values {
			if index == 0 {
				where += field + " like ?"
			} else {
				where += " or " + field + " like ?"
			}
			likes = append(likes, "%"+value+"%")
		}
	}
	where += ")"

	return db.Where(where, likes...), nil
}
