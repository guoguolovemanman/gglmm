package gglmm

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"strings"
)

// error
var (
	ErrConfigFile = errors.New("配置文件错误")
)

// ConfigString 字符配置
type ConfigString struct {
	Value string `json:"value"`
	Name  string `json:"name"`
}

// Status
var (
	StatusInvalid = ConfigString{Value: "invalid", Name: "无效"}
	StatusFrozen  = ConfigString{Value: "frozen", Name: "冻结"}
	StatusValid   = ConfigString{Value: "valid", Name: "有效"}
	Statuses      = []ConfigString{StatusValid, StatusFrozen, StatusInvalid}
)

// ConfigHTTP --
type ConfigHTTP struct {
	Address string
}

// Check --
func (config ConfigHTTP) Check() bool {
	if config.Address == "" || !strings.Contains(config.Address, ":") {
		return false
	}
	log.Println("ConfigHTTP check pass")
	return true
}

// ConfigRPC --
type ConfigRPC struct {
	Network string
	Address string
	Call    string
}

// Check --
func (config ConfigRPC) Check() bool {
	if config.Network == "" {
		return false
	}
	if config.Address == "" || !strings.Contains(config.Address, ":") {
		return false
	}
	if config.Call == "" {
		return false
	}
	log.Println("ConfigRPC check pass")
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
	log.Println("ConfigDB check pass")
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
