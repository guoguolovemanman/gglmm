package gglmm

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
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

// ConfigJWT --
type ConfigJWT struct {
	Expires int64
	Secret  string
}

// Check --
func (config *ConfigJWT) Check(cmd string) bool {
	if cmd == "all" || cmd == "write" {
		if config.Expires <= 0 {
			return false
		}
	}
	if cmd == "all" || cmd == "read" {
		if config.Secret == "" {
			return false
		}
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
		return errors.New("config check fail")
	}
	return nil
}
