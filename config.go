package gglmm

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/jinzhu/gorm"
)

// error
var (
	ErrConfigFile         = errors.New("配置文件错误")
	ErrRequest            = errors.New("请求参数错误")
	ErrFilter             = errors.New("过滤参数错误")
	ErrFilterValueType    = errors.New("过滤值类型错误")
	ErrFilterValueSize    = errors.New("过滤值大小错误")
	ErrFilterOperate      = errors.New("过滤操作错误")
	ErrAction             = errors.New("不支持Action")
	ErrModelType          = errors.New("模型类型错误")
	ErrModelCanNotDeleted = errors.New("模型不可删除")
	ErrPathVar            = errors.New("路径参数错误")

	ErrGormRecordNotFound = gorm.ErrRecordNotFound
)

// ConfigInt8 --
type ConfigInt8 struct {
	Value int8   `json:"value"`
	Name  string `json:"name"`
}

// ConfigString --
type ConfigString struct {
	Value string `json:"value"`
	Name  string `json:"name"`
}

// status
var (
	StatusInvalid = ConfigInt8{Value: -128, Name: "无效"}
	StatusFrozen  = ConfigInt8{Value: -127, Name: "冻结"}
	StatusValid   = ConfigInt8{Value: 1, Name: "有效"}
	Statuses      = []ConfigInt8{StatusValid, StatusFrozen, StatusInvalid}
)

// ConfigAPI --
type ConfigAPI struct {
	Address string
}

// Check --
func (config ConfigAPI) Check() bool {
	if config.Address == "" || !strings.Contains(config.Address, ":") {
		return false
	}
	return true
}

// ConfigRPC --
type ConfigRPC struct {
	Network string
	Address string
}

// Check --
func (config ConfigRPC) Check() bool {
	if config.Network == "" {
		return false
	}
	if config.Address == "" || !strings.Contains(config.Address, ":") {
		return false
	}
	return true
}

// ConfigDB --
type ConfigDB struct {
	Dialect         string
	Address         string
	User            string
	Password        string
	Database        string
	MaxOpen         int
	MaxIdel         int
	ConnMaxLifetime int
}

// Check --
func (config ConfigDB) Check() bool {
	if config.Dialect == "" {
		return false
	}
	if config.Address == "" || !strings.Contains(config.Address, ":") {
		return false
	}
	if config.User == "" || config.Password == "" {
		return false
	}
	if config.Database == "" {
		return false
	}
	if config.MaxOpen <= 0 || config.MaxIdel <= 0 || config.ConnMaxLifetime <= 0 {
		return false
	}
	return true
}

// ConfigChecker 配置检查
type ConfigChecker interface {
	Check() bool
}

// ParseConfigFile --
func ParseConfigFile(file string, config ConfigChecker) error {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, config)
	if err != nil {
		return err
	}
	if !config.Check() {
		return ErrConfigFile
	}
	return nil
}
